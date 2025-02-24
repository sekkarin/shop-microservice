pipeline {
    agent any
    environment {
        ZAP_IMAGE = 'zaproxy/zap-stable:2.16.0'  // Docker image for OWASP ZAP
        GOLANG_IMAGE = 'golang:1.23'
        TRIVY_IMAGE = 'aquasec/trivy:0.59.1'

        IMAGE_NAME = 'sekkarindev/shop-microservice'
        TARGET_URL = 'http://localhost:3000'  // URL of the app you want to scan
        ZAP_PORT = '80'  // Port that ZAP will use
        ZAP_WAIT_TIME = '30'  // Wait for ZAP container to initialize
    }

    stages {
        stage('Fetch Code') {
            steps {
                checkout scmGit(branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[url: 'https://github.com/sekkarin/shop-microservice.git']])
            }
        }
        stage('Unit Tests') {
            steps {
                script {
                    sh '''
                    docker run --rm -v $WORKSPACE:/app -w /app ${GOLANG_IMAGE} sh -c "go mod tidy &&
                        cd __test__ &&
                        go test ./... -v -coverprofile=coverage.out | tee go-test-results.txt"
                    '''
                }
            }
            // post {
            //     always {
            //         archiveArtifacts artifacts: 'go-test-results.txt', fingerprint: true
            //     }
            // }
        }
        stage('SAST - Code Security Scan') {
            environment {
                SCANNER_HOME = tool 'SonarScanner'  // Make sure 'SonarScanner' matches Jenkins tool name
                SONARQUBE_SERVER = 'sonarqube-server'      // Make sure 'SonarQube' matches configured server in Jenkins
            }

            steps {
                script {
                    withSonarQubeEnv('sonarqube-server') {
                        sh "${SCANNER_HOME}/bin/sonar-scanner"
                    }
                }
            }
        }
        stage('SCA - Dependency Scan') {
            steps {
                script {
                    sh '''
                    docker run --rm  -v $WORKSPACE:/app ${TRIVY_IMAGE} fs --format template --template "@contrib/html.tpl" -o /app/SCA-report.html --scanners vuln,misconfig,secret,license /app
                    '''
                }
            }
            post {
                success {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: false,
                        keepAll: false, reportDir: '/var/lib/jenkins/workspace/Shop-microservices',
                        reportFiles: 'SCA-report.html',
                        reportName: 'HTML Report SCA',
                        reportTitles: '',
                        useWrapperFileDirectly: true
                    ])
                }
            }
        }
        stage('Build & Container Security Scan') {
            steps {
                script {
                    sh '''
                     docker build -t ${IMAGE_NAME}:latest -t ${IMAGE_NAME}:$BUILD_NUMBER -f ./build/Dockerfile .
                    '''
                    sh '''
                    docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock -v $WORKSPACE:/app ${TRIVY_IMAGE} image --format template --template "@contrib/html.tpl" -o /app/CSS-report.html --scanners vuln,misconfig,secret,license ${IMAGE_NAME}:$BUILD_NUMBER
                    '''
                }
            }
            post {
                success {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: false,
                        keepAll: false, reportDir: '/var/lib/jenkins/workspace/Shop-microservices',
                        reportFiles: 'CSS-report.html.html',
                        reportName: 'HTML Report CSS',
                        reportTitles: '',
                        useWrapperFileDirectly: true
                    ])
                }
            }
        }
        stage('DAST - Web Security Scan') {
            steps {
                script {
                    withCredentials([
                        file(credentialsId: 'VAULT_SECRET_ID', variable: 'SECRET_ID'),
                        file(credentialsId: 'VAULT_SECRET_TOKEN', variable: 'SECRET_TOKEN')
                    ]) {
                        sh 'mv $SECRET_ID ./vault-agent-config/'
                        sh 'mv $SECRET_TOKEN ./vault-agent-config/'
                        sh 'docker compose -f compose.yaml up -d --build'
                        sh '''
                            docker run --rm --user root -v ${WORKSPACE}:/zap/wrk $ZAP_IMAGE zap-api-scan.py -t http://$(ip -f inet -o addr show docker0 | awk '{print $4}' | cut -d '/' -f 1):3000/auth_v1/auth/login -f openapi -I -r report-api.html -d
                        '''
                    }
                }
            }
            post {
                always {
                    sh 'docker compose -f compose.yaml down'
                    sh 'rm -r ./vault-agent-config'
                }
                success {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: false,
                        keepAll: false, reportDir: '/var/lib/jenkins/workspace/Shop-microservices',
                        reportFiles: 'report-api.html',
                        reportName: 'HTML Report',
                        reportTitles: '',
                        useWrapperFileDirectly: true
                    ])
                }
            }
        // stage('Deploy to Kubernetes') {
        //     steps {
        //         script {
        //             sh 'echo "Deploy..........."'
        //         }
        //     }
        // }
        }
    }
}
