pipeline {
    agent any

    environment {
        NAME = "Production-grade DevOps"
    }

    stages {

        stage('Checkout Code') {
            steps {
                git url: 'https://github.com/christalselvin/Production-Grade-DevOps-Platform-for-Microservices.git'
                echo "Git clone completed"
            }
        }

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
                    sh 'pip install -r requirements.txt'
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
                echo "Docker images built"
            }
        }

    }
}
