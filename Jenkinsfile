pipeline {
    agent any

    stages {
        stage('Fetch Code') {
            steps {
                checkout scmGit(branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[url: 'https://github.com/sekkarin/ConsoleApp-with-sonarQ-jenkins.git']])
            }
        }
      
        stage('Build & Test in Docker') {
            steps {
                script {
                    sh """
                    echo "build........"
                    """
                }
            }
        }

          stage('Code Analysis') {
            environment {
                SCANNER_HOME = tool 'SonarScanner'  // Make sure 'SonarScanner' matches Jenkins tool name
                SONARQUBE_SERVER = 'sonarqube-server'      // Make sure 'SonarQube' matches configured server in Jenkins
            }

            steps {
                script {
                    withSonarQubeEnv('sonarqube-server') {
                        sh "${SCANNER_HOME}/bin/sonar-scanner \
                            -Dsonar.projectKey=shop-mricoservices \
                            -Dsonar.sources=."
                    }
                }
            }
        }
        
        stage('Scan security') {
            steps {
                script {
                    sh '...'
                }
            }
        }
        stage('Deploy') {
            steps {
                script {
                    sh 'echo "Deploy..........."'
                }
            }
        }
    }
}