import { ethers, waffle } from "hardhat"
import { BigNumber } from "ethers"
import { Gravity } from "../typechain/Gravity"
import { TestNebula } from "../typechain/TestNebula"
import { expect } from "./shared/expect"
import { testNebulaFixture } from "./shared/fixtures"

const emptyBytes32: string = "0x0000000000000000000000000000000000000000000000000000000000000000"
const emptyAddress: string = "0x0000000000000000000000000000000000000000"

export interface Queue {
  first: string;
  last: string;
}
export interface Subscription {
  owner: string;
  contractAddress: string;
  minConfirmations: number;
  reward: BigNumber;
}
export interface Pulse {
  dataHash: string;
  height: BigNumber;
}

describe("Nebula", () => {
  const [wallet, other, consul1, consul2, consul3] = waffle.provider.getWallets()

  let gravity: Gravity
  let nebula: TestNebula

  let loadFixture: ReturnType<typeof waffle.createFixtureLoader>

  before("create fixture loader", async () => {
    loadFixture = waffle.createFixtureLoader([wallet, other, consul1, consul2, consul3])
  })

  beforeEach("deploy test contracts", async () => {
    ;({ gravity, nebula } = await loadFixture(testNebulaFixture))
  })

  it("constructor initializes variables", async () => {
    expect(await nebula.bftValue()).to.eq(3)
    expect(await nebula.dataType()).to.eq(2)
    expect(await nebula.oracles(0)).to.eq(consul1.address)
    expect(await nebula.oracles(1)).to.eq(consul2.address)
    expect(await nebula.oracles(2)).to.eq(consul3.address)
    expect(await nebula.gravityContract()).to.eq(gravity.address)
  })

  it("starting state after deployment", async () => {
    expect(await nebula.rounds(0)).to.eq(false)

    var oracleQueue = await nebula.oracleQueue() as Queue
    expect(oracleQueue.first).to.eq(emptyBytes32)
    expect(oracleQueue.last).to.eq(emptyBytes32)
    expect(await nebula.oracleNextElement(emptyBytes32)).to.eq(emptyBytes32)
    expect(await nebula.oraclePrevElement(emptyBytes32)).to.eq(emptyBytes32)

    var subscriptionsQueue = await nebula.subscriptionsQueue() as Queue
    expect(subscriptionsQueue.first).to.eq(emptyBytes32)
    expect(subscriptionsQueue.last).to.eq(emptyBytes32)
    expect(await nebula.subscriptionsNextElement(emptyBytes32)).to.eq(emptyBytes32)
    expect(await nebula.subscriptionsPrevElement(emptyBytes32)).to.eq(emptyBytes32)

    var pulseQueue = await nebula.pulseQueue() as Queue
    expect(pulseQueue.first).to.eq(emptyBytes32)
    expect(pulseQueue.last).to.eq(emptyBytes32)
    expect(await nebula.pulseNextElement(emptyBytes32)).to.eq(emptyBytes32)
    expect(await nebula.pulsePrevElement(emptyBytes32)).to.eq(emptyBytes32)

    await expect(nebula.subscriptionIds(0)).to.be.reverted

    expect(await nebula.lastPulseId()).to.eq(0)

    var subscription = await nebula.subscriptions(emptyBytes32) as Subscription
    expect(subscription.owner).to.eq(emptyAddress)
    expect(subscription.contractAddress).to.eq(emptyAddress)
    expect(subscription.minConfirmations).to.eq(0)
    expect(subscription.reward).to.eq("0")

    var pulse = await nebula.pulses(0) as Pulse
    expect(pulse.dataHash).to.eq(emptyBytes32)
    expect(pulse.height).to.eq("0")

    expect(await nebula.isPulseSubSent(0, emptyBytes32)).to.eq(false)
  })

  describe("#receive", () => {
    it("", async () => {
    })
  })

  describe("#getOracles", () => {
    it("", async () => {
    })
  })

  describe("#getSubscribersIds", () => {
    it("", async () => {
    })
  })

  describe("#hashNewOracles", () => {
    it("", async () => {
    })
  })

  describe("#sendHashValue", () => {
    it("", async () => {
    })
  })

  describe("#updateOracles", () => {
    it("", async () => {
    })
  })

  describe("#sendValueToSubByte", () => {
    it("", async () => {
    })
  })

  describe("#sendValueToSubInt", () => {
    it("", async () => {
    })
  })

  describe("#sendValueToSubString", () => {
    it("", async () => {
    })
  })

  describe("#subscribe", () => {
    it("", async () => {
    })
  })
})
