const Gravity = artifacts.require("./Gravity/Gravity.sol");
const Queue = artifacts.require("./libs/QueueLib");
const Nebula = artifacts.require("./Nebula/Nebula.sol");

module.exports = async function(deployer, network, accounts) {
  await deployer.deploy(Queue);
  await deployer.link(Queue, Nebula);
//  await deployer.link(Queue, Gravity);

  await deployer.deploy(Gravity, [ accounts[0],accounts[1]  ], 1);
 let nubula = await deployer.deploy(Nebula, 0, "0xa3a2d245E621b7B38F0C136409aB72fbEA5106b0", accounts, 1);
 // let sub = await deployer.deploy(SubMock, nubula.address, "100000000000000");
 // await nubula.subscribe(sub.address, 0, "100000000000000")
 // await web3.eth.sendTransaction({ from: accounts[0], to: sub.address, value:  "10000000000000000000" });
};