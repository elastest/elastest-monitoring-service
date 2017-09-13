node('docker') {
    stage "Container Prep"
        echo("The node is up")
        def mycontainer = docker.image('elastest/docker-siblings')
        mycontainer.pull()
        mycontainer.inside("-u jenkins -v /var/run/docker.sock:/var/run/docker.sock:rw -v ${WORKSPACE}:/home/jenkins/.m2") {
            git 'https://github.com/elastest/elastest-monitoring-service'

            stage "Tests"
                echo ("Starting tests")
                sh 'docker run -v $(pwd)/go_EMS:/go/go_EMS golang /bin/bash -c "cd go_EMS;go test"'

            stage "Build images - Package"
                echo ("Building full version")
                sh 'docker build -t elastest/ems:0.1 .'
                def myfullimage = docker.image('elastest/ems:0.1');
                echo ("Building min version")
                sh 'docker build -f Dockerfile_min -t elastest/ems_min:0.1 .'
                def myminimage = docker.image('elastest/ems_min:0.1');

            stage "Run images"
                myfullimage.run()
                myminimage.run()

            stage "Publish"
                echo ("Publishing")
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'elastestci-dockerhub',
                    usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
                    sh 'docker login -u "$USERNAME" -p "$PASSWORD"'
                    myfullimage.push()
					myminimage.run()
                }   
        }   
}
