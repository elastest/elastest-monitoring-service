node('docker') {
    stage "Container Prep"
        echo("The node is up")
        def mycontainer = docker.image('elastest/docker-in-docker:latest')
        mycontainer.pull()
        mycontainer.inside("-u jenkins -v /var/run/docker.sock:/var/run/docker.sock:rw") {
            git 'https://github.com/elastest/elastest-monitoring-service'

            stage "Tests"
                echo ("Starting tests")
                echo ("(No tests yet)")

            stage "Build image - Package"
                echo ("Building")
                sh 'docker build -t elastest/elastest-monitoring-service .'
                def myimage = docker.image('elastest/elastest-monitoring-service');

            stage "Run image"
                myimage.run()

            stage "Publish"
                echo ("Publishing")
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'elastestci-dockerhub',
                    usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
                    sh 'docker login -u "$USERNAME" -p "$PASSWORD"'
                    myimage.push()
                }   
        }   
}
