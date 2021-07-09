import { ethers, waffle } from "hardhat"
import { BigNumber } from "ethers"

import { Gravity } from "../../typechain/Gravity"
import { TestNebula } from "../../typechain/TestNebula"
import { SubMockBytes } from "../../typechain/SubMockBytes"
import { SubMockInt } from "../../typechain/SubMockInt"
import { SubMockString } from "../../typechain/SubMockString"

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
  subMockBytes: SubMockBytes
  subMockString: SubMockString
  subMockInt: SubMockInt
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

  const subMockBytesFactory = await ethers.getContractFactory("SubMockBytes")
  const subMockBytes = await subMockBytesFactory.deploy(nebula.address, 0)

  const subMockStringFactory = await ethers.getContractFactory("SubMockString")
  const subMockString = await subMockStringFactory.deploy(nebula.address, 0)

  const subMockIntFactory = await ethers.getContractFactory("SubMockInt")
  const subMockInt = await subMockIntFactory.deploy(nebula.address, 0)

  return { gravity, nebula, subMockBytes, subMockString, subMockInt }
}
