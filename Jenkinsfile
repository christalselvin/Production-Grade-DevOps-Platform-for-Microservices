pipeline {
    agent any

    environment {
        NAME = "Production-grade DevOps"
    }

    stages {

        stage('Build Order Service (Java)') {
            steps {
                dir('order-service') {
                    sh 'mvn clean package'
                    echo "Java done"
                }
            }
        }

        stage('Build Payment Service (Go)') {
            steps {
                dir('payment-service') {
                    sh 'go build -o payment-service'
                    echo "Go done"
                }
            }
        }

        stage('Build Product Service (Python)') {
    steps {
        dir('product-service') {
            sh '''
            python3 -m venv venv
            . venv/bin/activate
            pip install -r requirements.txt
            '''
            echo "Python done"
        }
    }
}

        stage('Build User Service (Node)') {
            steps {
                dir('user-service') {
                    sh 'npm install'
                    echo "Node done"
                }
            }
        }

        stage('Build Docker Images') {
            steps {
                sh 'docker-compose build'
            }
        }

    }
}
