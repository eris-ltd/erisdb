#! /bin/bash

echo "your new chain, kind marmot: $CHAIN_ID"

# lay the genesis
# if it exists, just overwrite the chain id
if [ ! -f $CHAIN_DIR/genesis.json ]; then
	if [ "$CSV" = "" ]; then
		mintgen random --dir="$CHAIN_DIR" 1 $CHAIN_ID
		ifExit "Error creating random genesis file"
	else
		mintgen known --csv="$CSV" $CHAIN_ID > $CHAIN_DIR/genesis.json
		ifExit "Error creating genesis file from csv"
	fi
else
	# apparently just outputing to $CHAIN_DIR/genesis.json doesn't work so we copy
	cat $CHAIN_DIR/genesis.json | jq .chain_id=\"$CHAIN_ID\" > genesis.json
	cp genesis.json $CHAIN_DIR/genesis.json
fi


# if no config was given, lay one with the given options
if [ ! -f $CHAIN_DIR/config.toml ]; then
	echo "running mintconfig $CONFIG_OPTS"
	mintconfig $CONFIG_OPTS > $CHAIN_DIR/config.toml
else
	echo "found config file:"
	cat $CHAIN_DIR/config.toml
fi

# run the node.
# TODO: maybe bring back this stopping option if we think its useful
# tendermint node & last_pid=$! && sleep 1 && kill -KILL $last_pid
if [ $ERISDB_API ]; then
	echo "Running chain $CHAIN_ID (via ErisDB API)"
	erisdb $CHAIN_DIR
	ifExit "Error starting erisdb"
else
	echo Running chain $CHAIN_ID
	tendermint node
	ifExit "Error starting tendermint"
fi
