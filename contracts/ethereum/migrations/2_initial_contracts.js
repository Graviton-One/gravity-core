const SubMock = artifacts.require("SubMock");
const Nebula = artifacts.require("Nebula");
const Queue = artifacts.require("./libs/QueueLib");

module.exports = async function(deployer, network, accounts) {
  await deployer.deploy(Queue);
  await deployer.link(Queue, Nebula);
  let nubula = await deployer.deploy(Nebula, ADDRESS, 1);
  let sub = await deployer.deploy(SubMock, nubula.address, "100000000000000");
  await nubula.subscribe(sub.address, 0, "100000000000000")
  await web3.eth.sendTransaction({ from: accounts[0], to: sub.address, value:  "10000000000000000000" });
};
