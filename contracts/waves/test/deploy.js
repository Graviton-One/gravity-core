const wvs = 10 ** 8;


let subMockSeed = "waves private node seed with waves tokens submock"
let nebulaPubKey = "2K3zsM6XaqxaedbuC6dRB8cVX8TcnGRAXSkRyUmXiSAj"
let nebulaSeed = "waves private node seed with waves tokens nebula"
let oracles = [
    "AYoAcvavnQtYWU9E62PUXh56vdJipN1Q7JZ49MGw3nTM",
    "7yvzEhNziyBxDVp6o85C6FJUp3DdhPhmF8CNB3GcgQ48",
    "GFDQU6mvL6dPRTRTFsGQd6UuPdXYekqv691fKLKaUMmJ",
    "EtjJzM4s2FnsPz3Vt3YPLoei1745sqH63gsoZKA2LBCx",
    "CJHr6t5jpcJCrCsutW7BBE87deUDdxcBmST5zdCgSJqg"
]
describe('Deploy script', async function () {
    it('Deploy contract', async function () {
        const setScriptNebulaTx = setScript({ script: compile(file("../script/nebula.ride")), fee: 1400000,}, nebulaSeed); 
        await broadcast(setScriptNebulaTx)

        const setScriptSubMockTx = setScript({ script: compile(file("../script/subMock.ride")), fee: 1400000,}, subMockSeed); 
        await broadcast(setScriptSubMockTx)

        const constructorData = data({
            data: [
                { key: "oracles", value: oracles[0] + "," + oracles[1] + "," + oracles[2] + "," + oracles[3] + "," + oracles[4]},
                { key: "bft_coefficient", value: 1 },
                { key: "subscriber_address", value: address(subMockSeed) },
                { key: "contract_pubkey", value: nebulaPubKey }
            ],
            fee: 500000
        }, nebulaSeed);
        await broadcast(constructorData)

        console.log("Nebula:" + address(nebulaSeed))
        console.log("SubMock:" + address(subMockSeed))
    })
})