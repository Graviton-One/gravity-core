const SubMock = artifacts.require("SubMock");
const Nebula = artifacts.require("Nebula");
const Queue = artifacts.require("./libs/QueueLib");

module.exports = async function(deployer, network, accounts) {
  await deployer.deploy(Queue);
  await deployer.link(Queue, Nebula);
  let nubula = await deployer.deploy(Nebula, "test#1", accounts);
  let sub = await deployer.deploy(SubMock, nubula.address, "100000000000000000");
  await web3.eth.sendTransaction({ from: accounts[0], to: sub.address, value:  "1000000000000000000" });
};
