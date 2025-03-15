/* groovylint-disable LineLength */
pipeline {
    agent any
    environment {
        ZAP_IMAGE = 'zaproxy/zap-stable:2.16.0'  // Docker image for OWASP ZAP
        GOLANG_IMAGE = 'golang:1.23'
        TRIVY_IMAGE = 'aquasec/trivy:0.59.1'

        IMAGE_NAME = 'shop-microservice'
        TARGET_URL = 'http://localhost:3000'  // URL of the app you want to scan
        ZAP_PORT = '80'  // Port that ZAP will use
        ZAP_WAIT_TIME = '30'  // Wait for ZAP container to initialize

        HARBOR_REGISTRY = 'harbor.warering.online'
        HARBOR_PROJECT =  'shop-microservices'
        NAME_IMAGE_WITH_REGISTY = "${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${IMAGE_NAME}"
        SECRETS_DIR = './secrets-prod'

        CHART_NAME = 'auth-service'           // Change to your Helm chart name
        CHART_VERSION = "1.0.${BUILD_NUMBER}"

        GIT_CREDENTIALS_ID = 'github-ssh'
        GIT_REPO_URL = 'git@github.com:sekkarin/shop-microservice.git'
        GIT_BRANCH = 'main'
    }

    stages {
        stage('Fetch Code') {
            steps {
                checkout scmGit(branches: [[name: GIT_BRANCH]],
                extensions: [ cleanBeforeCheckout(deleteUntrackedNestedRepositories: true),
                [$class: 'WipeWorkspace']],
                userRemoteConfigs: [[credentialsId: GIT_CREDENTIALS_ID, url: GIT_REPO_URL]])
                script {
                    def changes = sh(script: 'git diff --name-only HEAD~1', returnStdout: true).trim()
                    def servicesToDeploy = []
                    echo "Changed Files:\n${changes}"
                    if (changes.contains('modules/auth/') || changes.contains('server/auth.go')) {
                        servicesToDeploy << 'auth'
                    }
                    if (changes.contains('modules/inventory/') || changes.contains('server/inventory.go')) {
                        servicesToDeploy << 'inventory'
                    }
                    if (changes.contains('modules/item/') || changes.contains('server/item.go')) {
                        servicesToDeploy << 'item'
                    }
                    if (changes.contains('modules/payment/') || changes.contains('server/payment.go')) {
                        servicesToDeploy << 'payment'
                    }
                    if (changes.contains('modules/player/') || changes.contains('server/player.go')) {
                        servicesToDeploy << 'player'
                    }
                    env.SERVICES_TO_DEPLOY = servicesToDeploy.join(' ')
                }
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
        //     // post {
        //     //     always {
        //     //         archiveArtifacts artifacts: 'go-test-results.txt', fingerprint: true
        //     //     }
        //     // }
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
        //             docker run --rm  -v $WORKSPACE:/app ${TRIVY_IMAGE} fs --format template --template "@contrib/html.tpl" -o /app/SCA-report.html --scanners vuln,misconfig,secret,license /app
        //             '''
        //         }
        //     }
        //     post {
        //         success {
        //             publishHTML([
        //                 allowMissing: false,
        //                 alwaysLinkToLastBuild: false,
        //                 keepAll: false, reportDir: '/var/lib/jenkins/workspace/Shop-microservices',
        //                 reportFiles: 'SCA-report.html',
        //                 reportName: 'HTML Report SCA',
        //                 reportTitles: '',
        //                 useWrapperFileDirectly: true
        //             ])
        //         }
        //     }
        // }
        stage('Build & Container Security Scan') {
            steps {
                script {
                    sh '''
                     docker build -t ${NAME_IMAGE_WITH_REGISTY}:latest -t ${NAME_IMAGE_WITH_REGISTY}:$BUILD_NUMBER -f ./Dockerfile .
                    '''
                    sh '''
                    docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock -v $WORKSPACE:/app ${TRIVY_IMAGE} image --format template --template "@contrib/html.tpl" -o /app/CSS-report.html --scanners vuln,misconfig,secret,license ${NAME_IMAGE_WITH_REGISTY}:$BUILD_NUMBER
                    '''
                }
            }
            post {
                success {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: false,
                        keepAll: false, reportDir: '/var/lib/jenkins/workspace/Shop-microservices',
                        reportFiles: 'CSS-report.html',
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
                    // sh '''
                    //     docker run --rm --user root  -v ${WORKSPACE}:/zap/wrk $ZAP_IMAGE zap-api-scan.py -t http://$(ip -f inet -o addr show docker0 | awk '{print $4}' | cut -d '/' -f 1):3000/auth_v1/auth/login -f openapi -I -r report-api.html
                    // '''
                    // sh '''
                    //     docker run --rm --user root  -v ${WORKSPACE}:/zap/wrk $ZAP_IMAGE zap-api-scan.py -t http://$(ip -f inet -o addr show docker0 | awk '{print $4}' | cut -d '/' -f 1):3000/auth_v1/auth/refresh-token -f openapi -I -r report-api.html
                    // '''
                    // sh '''
                    //     docker run --rm --user root  -v ${WORKSPACE}:/zap/wrk $ZAP_IMAGE zap-api-scan.py -t http://$(ip -f inet -o addr show docker0 | awk '{print $4}' | cut -d '/' -f 1):3000/auth_v1/auth/logout -f openapi -I -r report-api.html
                    // '''
                    // sh '''
                    //     docker run --rm --user root  -v ${WORKSPACE}:/zap/wrk $ZAP_IMAGE zap-api-scan.py -t http://$(ip -f inet -o addr show docker0 | awk '{print $4}' | cut -d '/' -f 1):3000/auth_v1 -f openapi -I -r report-api.html
                    // '''
                    }
                    withCredentials([usernamePassword(credentialsId: 'JenkinsCredential', usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                        sh "docker login $HARBOR_REGISTRY -u $HARBOR_USER -p $HARBOR_PASS"
                        sh "docker push $NAME_IMAGE_WITH_REGISTY:latest"
                        sh "docker push $NAME_IMAGE_WITH_REGISTY:$BUILD_NUMBER"
                    }
                }
            }
            post {
                always {
                    sh 'docker compose -f compose.yaml down'
                    sh "docker rmi $NAME_IMAGE_WITH_REGISTY:$BUILD_NUMBER"
                    sh "docker rmi $NAME_IMAGE_WITH_REGISTY:latest"
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
        }
        stage('Push Helm Chart') {
            steps {
                script {
                    withCredentials([
                        file(credentialsId: 'VAULT_PROD_ENV_SECRET_ID', variable: 'SECRET_ID')
                    ]) {
                        sh 'mv $SECRET_ID ./vault-config/'
                        sh '''
                        if [ ! -d "secrets-prod" ]; then
                            mkdir secrets-prod
                        fi
                        '''
                        sh '''
                            docker run -d --rm \
                                --name vault-agent \
                                --entrypoint /bin/sh \
                                -e VAULT_ADDR=http://192.168.60.50:8200 \
                                -v ./vault-config:/etc/vault:rw \
                                -v ./secrets-prod:/vault/secrets:rw \
                                --cap-add IPC_LOCK \
                                --privileged \
                                --user $(id -u jenkins):$(id -g jenkins) \
                                hashicorp/vault:1.18 \
                                -c "mkdir -p /etc/vault && vault agent -config=/etc/vault/vault-agent.hcl"
                        '''
                        def subdirectoryCount = 0
                        while (subdirectoryCount < 6) {
                            subdirectoryCount = sh(script: "find ${SECRETS_DIR} -maxdepth 1 -type d | wc -l", returnStdout: true).trim().toInteger()
                            echo "Waiting for exactly 6 subdirectories in ${SECRETS_DIR}... (Current count: ${subdirectoryCount})"
                            sleep(time: 5, unit: 'SECONDS')
                        }
                        echo "Found 6 subdirectories in ${SECRETS_DIR}. Proceeding with the next step."
                        script {
                            script {
                                if (env.SERVICES_TO_DEPLOY?.trim()) {  // Check if SERVICES_TO_DEPLOY is not empty
                                    def services = env.SERVICES_TO_DEPLOY.split(' ')

                                    withCredentials([usernamePassword(credentialsId: 'JenkinsCredential', usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                                        sh "helm registry login ${HARBOR_REGISTRY} --username ${HARBOR_USER} --password ${HARBOR_PASS}"
                                        for (service in services) {
                                            if (service.trim()) {  // Ensure no empty values
                                                sh "cp -r ./secrets-prod/${service}-prod/secret.yaml ./charts/${service}/${service}-service/templates/secret.yaml"
                                                sh "helm dependency update ./charts/${service}/${service}-service/"
                                                sh "helm package ./charts/${service}/${service}-service --version ${CHART_VERSION}"
                                                sh "helm push ${service}-service-${CHART_VERSION}.tgz oci://${HARBOR_REGISTRY}/${HARBOR_PROJECT}"
                                                sh "rm -rf ${service}-service-${CHART_VERSION}.tgz"
                                                dir("applicationset/cluster-config/${service}-service") {
                                                    sh """
                                                  jq '.cluster.version = "${CHART_VERSION}"' config.json > temp.json && mv temp.json config.json
                                                 """
                                                }
                                            }
                                        }
                                        gitPush()
                                    }
                            } else {
                                    echo 'No services to deploy. Skipping deployment step.'
                                }
                            }
                        }
                    }
                }
            }
            post {
                always {
                    sh 'docker stop vault-agent'
                // sh 'rm -r ${SECRETS_DIR}'
                }
            }
        }
    }
    post {
        always {
            cleanWs(
            cleanWhenNotBuilt: false, // Don't clean if build wasn't executed
            deleteDirs: true,         // Delete all directories
            disableDeferredWipeout: true,  // Clean immediately after build
            notFailBuild: true,       // Ensure build doesn't fail due to cleanup
            patterns: [
                    [pattern: '*.properties', type: 'INCLUDE'], // Keep all .properties files
                    [pattern: '*/', type: 'EXCLUDE'] // Delete all directories
                ]
        )
        }
    }
}
