template='
{
    "ride_directory": "script",
    "test_directory": "test",
    "envs": {
        "custom": {
            "API_BASE": "'$1'",
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

printf "$template" > $2
