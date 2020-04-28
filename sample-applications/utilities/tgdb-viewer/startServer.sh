export GRAPH_BUILDER_HOME=$(pwd)

java -cp ./bin:./exeJar/tgdb-terminal-0.9.0.jar com.tibco.tge.main.TGDBViewer -s ./server.json
