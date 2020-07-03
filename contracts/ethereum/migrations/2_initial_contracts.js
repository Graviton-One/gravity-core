const SubMock = artifacts.require("SubMock");
const Nebula = artifacts.require("./Nebula/Nebula.sol");
const Gravity = artifacts.require("./Gravity/Gravity.sol");
const Queue = artifacts.require("./libs/QueueLib");

module.exports = async function(deployer, network, accounts) {
  await deployer.deploy(Queue);
  await deployer.link(Queue, Nebula);
  await deployer.link(Queue, Gravity);

  let gravity = await deployer.deploy(Gravity, [ accounts[0] ], 1);
  let nubula = await deployer.deploy(Nebula, accounts, gravity.address, 1);
  let sub = await deployer.deploy(SubMock, nubula.address, "100000000000000");
  await nubula.subscribe(sub.address, 0, "100000000000000")
  await web3.eth.sendTransaction({ from: accounts[0], to: sub.address, value:  "10000000000000000000" });
};
