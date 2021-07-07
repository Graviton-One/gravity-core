import { ethers, waffle } from "hardhat"
import { BigNumber } from "ethers"

import { Gravity } from "../../typechain/Gravity"
import { TestNebula } from "../../typechain/TestNebula"

import { Fixture } from "ethereum-waffle"

interface GravityFixture {
  gravity: Gravity
}

export const gravityFixture: Fixture<GravityFixture> =
  async function (
    [wallet, other, consul1, consul2, consul3],
    provider
  ): Promise<GravityFixture> {
  const gravityFactory = await ethers.getContractFactory("Gravity")
  const gravity = (await gravityFactory.deploy(
      [consul1.address,
       consul2.address,
       consul3.address],
      3
  )) as Gravity
  return { gravity }
}

interface TestNebulaFixture extends GravityFixture {
  nebula: TestNebula
}

export const testNebulaFixture: Fixture<TestNebulaFixture> =
  async function (
    [wallet, other, consul1, consul2, consul3],
    provider
  ): Promise<TestNebulaFixture> {
  const { gravity } = await gravityFixture(
    [wallet, other, consul1, consul2, consul3],
    provider
  )

  const queueFactory = await ethers.getContractFactory("QueueLib")
  const queue = await queueFactory.deploy()

  const nebulaFactory = await ethers.getContractFactory("TestNebula", {
    libraries: {
      QueueLib: queue.address,
    },
  })
  const nebula = (await nebulaFactory.deploy(
      2,
      gravity.address,
      [consul1.address,
       consul2.address,
       consul3.address],
      3
  )) as TestNebula
  return { gravity, nebula }
}
