pipeline {
    agent any

    options {
        ansiColor('xterm')
        timeout(time: 1, unit: 'HOURS')
        timestamps()
    }

    stages {
        stage('System Verification') {
            steps {
                script {
                    echo "CI/CD Environment Check"
                    echo "Running on: ${env.NODE_NAME}"
                    echo "Build Number: ${env.BUILD_NUMBER}"
                }
            }
        }

        stage('Docker Integration Test') {
            steps {
                echo 'Checking Docker CLI availability...'
                sh 'docker --version'
                
                echo 'Checking access to host Docker socket...'
                sh 'docker ps'
            }
        }

        stage('Environment Summary') {
            steps {
                sh '''
                    echo "Environment verification:"
                    echo "Java Version: $(java -version 2>&1 | head -n 1)"
                    echo "OS: $(cat /etc/os-release | grep PRETTY_NAME)"
                '''
            }
        }
    }

    post {
        success {
            echo 'Verification Pipeline Completed Successfully'
        }
        failure {
            echo 'Verification Pipeline Failed - Check Docker socket permissions'
        }
    }
}
