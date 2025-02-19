pipeline {
    agent any

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
                    echo "Unit Tests"
                    // if [ -f "composer.json" ]; then
                    //     composer install --no-interaction
                    //     vendor/bin/phpunit --log-junit unit-tests.xml
                    // elif [ -f "package.json" ]; then
                    //     npm install
                    //     npm test -- --ci --json --outputFile=unit-tests.json
                    // else
                    //     echo "No test framework detected. Skipping..."
                    // fi
                    '''
                }
            }
            // post {
            //     always {
            //         junit '**/unit-tests.xml'
            //         archiveArtifacts artifacts: '**/unit-tests.json', fingerprint: true
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
                    sh 'echo "Scan security"'
                }
            }
        }
        stage('Build & Container Security Scan') {
            steps {
                script {
                    sh 'echo "Scan security"'
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
