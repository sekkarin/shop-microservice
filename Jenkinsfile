pipeline {
    agent any
    environment {
        IMAGE_NAME = 'sekkarindev/shop-microservice'
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
                    docker run --rm -v $WORKSPACE:/app -w /app golang:1.23 sh -c "go mod tidy &&
                        cd __test__ &&
                        go test ./... -v -coverprofile=coverage.out | tee go-test-results.txt"
                    '''
                }
            }
            post {
                always {
                    archiveArtifacts artifacts: 'go-test-results.txt', fingerprint: true
                }
            }
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
                    docker run --rm  -v $WORKSPACE:/app aquasec/trivy:latest fs -f json -o /app/trivy-report.json --scanners vuln,misconfig,secret,license /app
                    '''
                }
            }
            post {
                always {
                    // Archive the Trivy report after the scan
                    archiveArtifacts artifacts: 'trivy-report.json', allowEmptyArchive: true
                }
            }
        }
        stage('Build & Container Security Scan') {
            steps {
                def commitId = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                script {
                    sh '''
                     docker build -t ${IMAGE_NAME}:latest -t ${IMAGE_NAME}:${commitId} .
                     docker images |grep sekkarindev/shop-microservice
                    '''
                }
            }
        }
        stage('DAST - Web Security Scan') {
            steps {
                script {
                    sh 'echo "Scan security"'
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
