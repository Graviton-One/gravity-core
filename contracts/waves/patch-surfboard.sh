
#!/bin/bash
api=''
out=''

while [ -n "$1" ]
do
	 case "$1" in
	   --api) api=$2 ;;
	   --out) out=$2 ;;
	 esac
	 shift
done


template='
{
    "ride_directory": "script",
    "test_directory": "test",
    "envs": {
        "custom": {
            "API_BASE": "%s",
            "CHAIN_ID": "R",
            "SEED": "waves private node seed with waves tokens",
            "timeout": 60000
        }
    },
    "defaultEnv": "custom",
    "mocha": {
        "timeout": 60000
    }
}'

printf "$template" "$api" > "$out"