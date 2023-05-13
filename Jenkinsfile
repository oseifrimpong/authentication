pipeline {
    agent any
    environment {
        AWS_ACCOUNT_ID="119796859579"
        AWS_DEFAULT_REGION="us-east-2"
        IMAGE_REPO_NAME="auth"
        IMAGE_TAG="develop"
        THE_BUTLER_SAYS_SO=credentials('jenkins-keyspecs')
        REPOSITORY_URI = "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/${IMAGE_REPO_NAME}"
    }

    stages {

         stage('Logging into AWS ECR') {
         agent {label 'linux-1'}
            steps {
                script {
                sh "aws ecr get-login-password --region ${AWS_DEFAULT_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com"
                }

            }
        }
    // Building Docker images
    stage('Building image') {
    agent {label 'linux-1'}
      steps{
        script {
          sh "docker build -f deploy/docker/Dockerfile -t ${IMAGE_REPO_NAME}:${IMAGE_TAG} . --network host"
        }
      }
    }

    // Uploading Docker images into AWS ECR
    stage('Pushing to ECR') {
    agent {label 'linux-1'}
     steps{
         script {
                sh "docker tag ${IMAGE_REPO_NAME}:${IMAGE_TAG} ${REPOSITORY_URI}:$IMAGE_TAG"
                sh "docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/${IMAGE_REPO_NAME}:${IMAGE_TAG}"
         }
        }
      }
    stage('Deployment on  to dev ') { 
    agent {label 'regilex-dev'}  
  when {
    branch "PR-*"
      }
      
     steps{
         script {
                sh "docker service update regilex_auth"
                
         }
        }
      }
    }
}
