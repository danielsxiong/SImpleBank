pipeline {
  agent any

  stages {
    stage('Dev') {
      steps {
        echo 'Checkout from git...'
        git 'https://github.com/danielsxiong/SImpleBank.git'
        echo 'Building'
        make build
      }
    }
  }
}