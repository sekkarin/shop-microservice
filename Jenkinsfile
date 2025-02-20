pipeline {
    agent any
    environment {
        IMAGE_NAME = 'sekkarindev/shop-microservice'
        ZAP_IMAGE = 'zaproxy/zap-stable:2.16.0'  // Docker image for OWASP ZAP
        TARGET_URL = 'http://your-app-url.com' // Replace with your application's URL
        GOLANG_IMAGE = 'golang:1.23'
        TRIVY_IMAGE = 'aquasec/trivy:0.59.1'
    }

    stages {
        stage('Fetch Code') {
            steps {
                checkout scmGit(branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[url: 'https://github.com/sekkarin/shop-microservice.git']])
            }
        }
        // stage('Unit Tests') {
        //     steps {
        //         script {
        //             sh '''
        //             docker run --rm -v $WORKSPACE:/app -w /app ${GOLANG_IMAGE} sh -c "go mod tidy &&
        //                 cd __test__ &&
        //                 go test ./... -v -coverprofile=coverage.out | tee go-test-results.txt"
        //             '''
        //         }
        //     }
        //     post {
        //         always {
        //             archiveArtifacts artifacts: 'go-test-results.txt', fingerprint: true
        //         }
        //     }
        // }
        // stage('SAST - Code Security Scan') {
        //     environment {
        //         SCANNER_HOME = tool 'SonarScanner'  // Make sure 'SonarScanner' matches Jenkins tool name
        //         SONARQUBE_SERVER = 'sonarqube-server'      // Make sure 'SonarQube' matches configured server in Jenkins
        //     }

        //     steps {
        //         script {
        //             withSonarQubeEnv('sonarqube-server') {
        //                 sh "${SCANNER_HOME}/bin/sonar-scanner"
        //             }
        //         }
        //     }
        // }
        // stage('SCA - Dependency Scan') {
        //     steps {
        //         script {
        //             sh '''
        //             docker run --rm  -v $WORKSPACE:/app ${TRIVY_IMAGE} fs -f json -o /app/trivy-report.json --scanners vuln,misconfig,secret,license /app
        //             '''
        //         }
        //     }
        //     post {
        //         always {
        //             // Archive the Trivy report after the scan
        //             archiveArtifacts artifacts: 'trivy-report.json', allowEmptyArchive: true
        //         }
        //     }
        // }
        // stage('Build & Container Security Scan') {
        //     steps {
        //         script {
        //             sh '''
        //              docker build -t ${IMAGE_NAME}:latest -t ${IMAGE_NAME}:$BUILD_NUMBER -f ./build/Dockerfile .
        //             '''
        //             sh '''
        //             docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock -v $WORKSPACE:/app ${TRIVY_IMAGE} image -f json -o /app/trivy-report-image.json --scanners vuln,misconfig,secret,license ${IMAGE_NAME}:$BUILD_NUMBER
        //             '''
        //         }
        //     }
        //     post {
        //         always {
        //             archiveArtifacts artifacts: 'trivy-report-image.json', allowEmptyArchive: true
        //         }
        //     }
        // }
        stage('DAST - Web Security Scan') {
            steps {
                script {
                    sh '''
                        docker compose -f ./build/compose.yaml up -d
                    '''
                }
            }
        }
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    sh 'echo "Deploy..........."'
                }
            }
        }
    }
    post {
        failure {
            script {
                echo 'Security scan failed! Fix issues before proceeding.'
            }
        }
        success {
            script {
                echo 'Pipeline passed successfully!'
            }
        }
    }
}
