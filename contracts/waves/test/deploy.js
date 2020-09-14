const wvs = 10 ** 8;

let gravitySeed = "waves private node seed with waves tokens8"

let nebulaSeed = "waves private node seed with waves tokens nebula2"
let subMockSeed = "waves private node seed with waves tokens submock"
let nebulaPubKey = "2K3zsM6XaqxaedbuC6dRB8cVX8TcnGRAXSkRyUmXiSAj"
let consul = [
    "AYoAcvavnQtYWU9E62PUXh56vdJipN1Q7JZ49MGw3nTM",
    "AYoAcvavnQtYWU9E62PUXh56vdJipN1Q7JZ49MGw3nTM",
    "AYoAcvavnQtYWU9E62PUXh56vdJipN1Q7JZ49MGw3nTM",
    "AYoAcvavnQtYWU9E62PUXh56vdJipN1Q7JZ49MGw3nTM",
    "AYoAcvavnQtYWU9E62PUXh56vdJipN1Q7JZ49MGw3nTM"
]

let main = "chapter ill milk glance mesh shoulder vivid derive fiber able tide trust diet violin exile"

describe('Deploy script', async function () {
    it('Deploy contract', async function () {
        let signedTranserTx = transfer({
            amount: 100000000,
            recipient: address(gravitySeed)
        }, main)
        await broadcast(signedTranserTx)

        signedTranserTx = transfer({
            amount: 100000000,
            recipient: address(nebulaSeed)
        }, main)
        await broadcast(signedTranserTx)

        /*  signedTranserTx = transfer({
                amount: 50000000,
                recipient: address(nebulaSeed)
            }, main)
            await broadcast(signedTranserTx)*/
     /*   signedTranserTx = transfer({
            amount: 10000000000,
            recipient: address("waves private node seed with waves tokens1")
        }, main)
        await broadcast(signedTranserTx)
*/

        const setScriptGravityTx = setScript({ script: compile(file("../script/gravity.ride")), fee: 1400000, }, gravitySeed);
        await broadcast(setScriptGravityTx)

        const setScriptNebulaTx = setScript({ script: compile(file("../script/nebula.ride")), fee: 1400000, }, nebulaSeed);
        await broadcast(setScriptNebulaTx)
  /*
        const setScriptSubMockTx = setScript({ script: compile(file("../script/subMock.ride")), fee: 1400000, }, subMockSeed);
        await broadcast(setScriptSubMockTx)*/

        let constructorData = data({
            data: [
                { key: "consuls_0", value: consul[0] + "," + consul[1] + "," + consul[2] + "," + consul[3] + "," + consul[4] },
                { key: "bft_coefficient", value: 1 },
                { key: "last_round", value: 0 },
            ],
            fee: 500000
        }, gravitySeed);
        await broadcast(constructorData)

        constructorData = data({
            data: [
                { key: "type", value: 0 },
                { key: "bft_coefficient", value: 1 },
                { key: "subscriber_address", value: address(subMockSeed) },
                { key: "gravity_contract", value: address(gravitySeed) },
                { key: "contract_pubkey", value: nebulaPubKey },
                { key: "oracles", value: consul[0] + "," + consul[1] + "," + consul[2] + "," + consul[3] + "," + consul[4] },
            ],
            fee: 500000
        }, nebulaSeed);
        await broadcast(constructorData)

        console.log("Gravity:" + address(gravitySeed))
        console.log("Nebula:" + address(nebulaSeed))
        console.log("SubMock:" + address(subMockSeed))
    })
})