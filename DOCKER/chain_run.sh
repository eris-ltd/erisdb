#! /bin/bash

if [ $ERISDB_API ]; then
	echo "Running chain $CHAIN_ID (via ErisDB API)"
	erisdb $CHAIN_DIR
	ifExit "Error starting erisdb"
else
	echo Running chain $CHAIN_ID
	tendermint node
	ifExit "Error starting tendermint"
fi
