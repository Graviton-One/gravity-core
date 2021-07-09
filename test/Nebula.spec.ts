import { ethers, waffle, network } from "hardhat"
import { BigNumber, ContractTransaction } from "ethers"
import { Gravity } from "../typechain/Gravity"
import { TestNebula } from "../typechain/TestNebula"
import { SubMockBytes } from "../typechain/SubMockBytes"
import { SubMockInt } from "../typechain/SubMockInt"
import { SubMockString } from "../typechain/SubMockString"
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
  const [wallet, other, consul1, consul2, consul3, oracle1, oracle2, oracle3] = waffle.provider.getWallets()

  let gravity: Gravity
  let nebula: TestNebula
  let subMockBytes: SubMockBytes
  let subMockString: SubMockString
  let subMockInt: SubMockInt

  let loadFixture: ReturnType<typeof waffle.createFixtureLoader>

  before("create fixture loader", async () => {
    loadFixture = waffle.createFixtureLoader([wallet, other, consul1, consul2, consul3])
  })

  beforeEach("deploy test contracts", async () => {
    ;({ gravity,
        nebula,
        subMockBytes,
        subMockString,
        subMockInt
      } = await loadFixture(testNebulaFixture))
  })

  function packAddresses(addresses: string[]): string {
    var pack: string = "0x"
    for (var i in addresses) {
      pack = ethers.utils.solidityPack([ "bytes", "address" ], [ pack, addresses[i] ]);
    }
    return pack
  }

  function hashAddresses(addresses: string[]): string {
    let pack = packAddresses(addresses)
    return ethers.utils.solidityKeccak256([ "bytes" ], [ pack ])
  }

  async function sendHashValue(hash: string): Promise<ContractTransaction> {

      const key1 = new ethers.utils.SigningKey(consul1.privateKey)
      const signature1 = await key1.signDigest(ethers.utils.arrayify(hash))
      let sig1 = ethers.utils.splitSignature(signature1);

      const key2 = new ethers.utils.SigningKey(consul2.privateKey)
      const signature2 = await key2.signDigest(ethers.utils.arrayify(hash))
      let sig2 = ethers.utils.splitSignature(signature2);

      const key3 = new ethers.utils.SigningKey(consul3.privateKey)
      const signature3 = await key3.signDigest(ethers.utils.arrayify(hash))
      let sig3 = ethers.utils.splitSignature(signature3);

      let vs = [sig1.v, sig2.v, sig3.v]
      let rs = [sig1.r, sig2.r, sig3.r]
      let ss = [sig1.s, sig2.s, sig3.s]

      return await nebula.sendHashValue(hash, vs, rs, ss)
  }

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

  // TODO:
  describe("#receive", () => {
    it("", async () => {
    })
  })

  describe("#getOracles", () => {
    it("returns array of oracles", async () => {
      let oracles = await nebula.getOracles()
      expect(oracles[0]).to.eq(consul1.address)
      expect(oracles[1]).to.eq(consul2.address)
      expect(oracles[2]).to.eq(consul3.address)
    })
  })

  describe("#getSubscribersIds", () => {
    it("returns empty array when there are no subscribers", async () => {
      let subscribers = await nebula.getSubscribersIds()
      expect(subscribers.length).to.eq(0)
    })

    it("returns array of subscriber ids", async () => {
    })
  })

  describe("#hashNewOracles", () => {
    it("hashes one address", async () => {
      let pack = ethers.utils.solidityPack([ "address" ], [ consul2.address ]);
      let hash = ethers.utils.solidityKeccak256([ "bytes" ], [ pack ])

      let hashNewOracles = await nebula.hashNewOracles([consul2.address])

      expect(hashNewOracles).to.eq(hash)
    })

    it("hashes three addresses", async () => {
      let pack1 = ethers.utils.solidityPack([ "address" ], [ consul1.address ]);
      let pack2 = ethers.utils.solidityPack([ "bytes", "address" ], [ pack1, consul2.address ]);
      let pack3 = ethers.utils.solidityPack([ "bytes", "address" ], [ pack2, consul3.address ]);
      let hash = ethers.utils.solidityKeccak256([ "bytes" ], [ pack3 ])

      let hashNewOracles = await nebula.hashNewOracles([consul1.address, consul2.address, consul3.address])

      expect(hashNewOracles).to.eq(hash)
    })

    it("hashes three addresses", async () => {
      let addresses = [consul1.address, consul2.address, consul3.address]
      let pack = packAddresses(addresses)
      let hash = ethers.utils.solidityKeccak256([ "bytes" ], [ pack ])

      let hashNewOracles = await nebula.hashNewOracles(addresses)

      expect(hashNewOracles).to.eq(hash)
    })

    it("hashes three addresses", async () => {
      let addresses = [consul1.address, consul2.address, consul3.address]
      let hash = hashAddresses(addresses)

      let hashNewOracles = await nebula.hashNewOracles(addresses)

      expect(hashNewOracles).to.eq(hash)
    })
  })

  describe("#sendHashValue", () => {
    it("fails when data hash is not signed by at least bft number of oracles", async () => {
      let hash = ethers.utils.solidityKeccak256(["bytes"], ["0x01"])

      const key1 = new ethers.utils.SigningKey(consul1.privateKey)
      const signature1 = await key1.signDigest(ethers.utils.arrayify(hash))
      let sig1 = ethers.utils.splitSignature(signature1);

      const key2 = new ethers.utils.SigningKey(consul2.privateKey)
      const signature2 = await key2.signDigest(ethers.utils.arrayify(hash))
      let sig2 = ethers.utils.splitSignature(signature2);

      const key3 = new ethers.utils.SigningKey(other.privateKey)
      const signature3 = await key3.signDigest(ethers.utils.arrayify(hash))
      let sig3 = ethers.utils.splitSignature(signature3);

      let vs = [sig1.v, sig2.v, sig3.v]
      let rs = [sig1.r, sig2.r, sig3.r]
      let ss = [sig1.s, sig2.s, sig3.s]

      await expect(nebula.sendHashValue(hash, vs, rs, ss))
        .to.be.revertedWith("invalid bft count")
    })

    it("updates data hash for the next pulse", async () => {
      let hash = ethers.utils.solidityKeccak256(["bytes"], ["0x01"])

      const key1 = new ethers.utils.SigningKey(consul1.privateKey)
      const signature1 = await key1.signDigest(ethers.utils.arrayify(hash))
      let sig1 = ethers.utils.splitSignature(signature1);

      const key2 = new ethers.utils.SigningKey(consul2.privateKey)
      const signature2 = await key2.signDigest(ethers.utils.arrayify(hash))
      let sig2 = ethers.utils.splitSignature(signature2);

      const key3 = new ethers.utils.SigningKey(consul3.privateKey)
      const signature3 = await key3.signDigest(ethers.utils.arrayify(hash))
      let sig3 = ethers.utils.splitSignature(signature3);

      let vs = [sig1.v, sig2.v, sig3.v]
      let rs = [sig1.r, sig2.r, sig3.r]
      let ss = [sig1.s, sig2.s, sig3.s]

      let tx = await nebula.sendHashValue(hash, vs, rs, ss)
      let rc = await tx.wait()
      let blockNumber = rc.blockNumber

      var pulse = await nebula.pulses(1) as Pulse
      expect(pulse.dataHash).to.eq(hash)
      expect(pulse.height).to.eq(blockNumber)
    })

    it("updates data hash for the next pulse", async () => {
      let hash = ethers.utils.solidityKeccak256(["bytes"], ["0x01"])

      let tx = await sendHashValue(hash)
      let rc = await tx.wait()
      let blockNumber = rc.blockNumber

      var pulse = await nebula.pulses(1) as Pulse
      expect(pulse.dataHash).to.eq(hash)
      expect(pulse.height).to.eq(blockNumber)
    })

    it("emits event", async () => {
      let hash = ethers.utils.solidityKeccak256(["bytes"], ["0x01"])

      const key1 = new ethers.utils.SigningKey(consul1.privateKey)
      const signature1 = await key1.signDigest(ethers.utils.arrayify(hash))
      let sig1 = ethers.utils.splitSignature(signature1);

      const key2 = new ethers.utils.SigningKey(consul2.privateKey)
      const signature2 = await key2.signDigest(ethers.utils.arrayify(hash))
      let sig2 = ethers.utils.splitSignature(signature2);

      const key3 = new ethers.utils.SigningKey(consul3.privateKey)
      const signature3 = await key3.signDigest(ethers.utils.arrayify(hash))
      let sig3 = ethers.utils.splitSignature(signature3);

      let vs = [sig1.v, sig2.v, sig3.v]
      let rs = [sig1.r, sig2.r, sig3.r]
      let ss = [sig1.s, sig2.s, sig3.s]

      let blockNumber = await wallet.provider.getBlockNumber()
      await expect(nebula.sendHashValue(hash, vs, rs, ss))
        .to.emit(nebula, "NewPulse")
        .withArgs(1, blockNumber+1, hash)
    })
  })

  describe("#updateOracles", () => {
    it("fails if new oracles are not signed by at least bft number oracles", async () => {
      let roundId = 1
      let addresses = [oracle1.address, oracle2.address, oracle3.address]
      let hash = hashAddresses(addresses)

      const key1 = new ethers.utils.SigningKey(consul1.privateKey)
      const signature1 = await key1.signDigest(ethers.utils.arrayify(hash))
      let sig1 = ethers.utils.splitSignature(signature1);

      const key2 = new ethers.utils.SigningKey(consul2.privateKey)
      const signature2 = await key2.signDigest(ethers.utils.arrayify(hash))
      let sig2 = ethers.utils.splitSignature(signature2);

      const key3 = new ethers.utils.SigningKey(other.privateKey)
      const signature3 = await key3.signDigest(ethers.utils.arrayify(hash))
      let sig3 = ethers.utils.splitSignature(signature3);

      let vs = [sig1.v, sig2.v, sig3.v]
      let rs = [sig1.r, sig2.r, sig3.r]
      let ss = [sig1.s, sig2.s, sig3.s]

      await expect(nebula.updateOracles(addresses, vs, rs, ss, roundId))
        .to.be.revertedWith("invalid bft count");
    })

    it("updates oracles", async () => {
      let roundId = 1
      let addresses = [oracle1.address, oracle2.address, oracle3.address]
      let hash = hashAddresses(addresses)

      const key1 = new ethers.utils.SigningKey(consul1.privateKey)
      const signature1 = await key1.signDigest(ethers.utils.arrayify(hash))
      let sig1 = ethers.utils.splitSignature(signature1);

      const key2 = new ethers.utils.SigningKey(consul2.privateKey)
      const signature2 = await key2.signDigest(ethers.utils.arrayify(hash))
      let sig2 = ethers.utils.splitSignature(signature2);

      const key3 = new ethers.utils.SigningKey(consul3.privateKey)
      const signature3 = await key3.signDigest(ethers.utils.arrayify(hash))
      let sig3 = ethers.utils.splitSignature(signature3);

      let vs = [sig1.v, sig2.v, sig3.v]
      let rs = [sig1.r, sig2.r, sig3.r]
      let ss = [sig1.s, sig2.s, sig3.s]

      await nebula.updateOracles(addresses, vs, rs, ss, roundId);

      let oracles = await nebula.getOracles()
      expect(oracles.length).to.eq(3)
      expect(oracles[0]).to.eq(oracle1.address)
      expect(oracles[1]).to.eq(oracle2.address)
      expect(oracles[2]).to.eq(oracle3.address)

      expect(await nebula.rounds(1)).to.eq(true)
    })
  })

  describe("#subscribe", () => {
    it("fails if a subscriber with matching parameters already exists", async () => {
      await nebula.subscribe(subMockBytes.address, 1, 0)
      await expect(nebula.subscribe(subMockBytes.address, 1, 0)).to.be.reverted
    })

    it("updates subscriptions", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      let subscription = await nebula.subscriptions(id) as Subscription
      expect(subscription.owner).to.eq(wallet.address)
      expect(subscription.contractAddress).to.eq(subMockBytes.address)
      expect(subscription.minConfirmations).to.eq(1)
      expect(subscription.reward).to.eq("0")
    })

    it("updates subscriptions queue with one subscriber", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      var subscriptionsQueue = await nebula.subscriptionsQueue() as Queue
      expect(subscriptionsQueue.first).to.eq(id)
      expect(subscriptionsQueue.last).to.eq(id)
      expect(await nebula.subscriptionsPrevElement(emptyBytes32)).to.eq(emptyBytes32)
      expect(await nebula.subscriptionsNextElement(emptyBytes32)).to.eq(emptyBytes32)
      expect(await nebula.subscriptionsPrevElement(emptyBytes32)).to.eq(emptyBytes32)
      expect(await nebula.subscriptionsNextElement(emptyBytes32)).to.eq(emptyBytes32)
    })

    it("updates subscriptions queue with two subscribers", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1_1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack1_2 = ethers.utils.solidityPack(["bytes"], [pack1_1])
      let id1 = ethers.utils.solidityKeccak256(["bytes"], [pack1_2])

      let pack2_1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockString.address, minConfirmations])
      let pack2_2 = ethers.utils.solidityPack(["bytes"], [pack2_1])
      let id2 = ethers.utils.solidityKeccak256(["bytes"], [pack2_2])

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      await nebula.subscribe(subMockString.address, minConfirmations, reward)

      var subscriptionsQueue = await nebula.subscriptionsQueue() as Queue
      expect(subscriptionsQueue.first).to.eq(id1)
      expect(subscriptionsQueue.last).to.eq(id2)
      expect(await nebula.subscriptionsPrevElement(emptyBytes32)).to.eq(emptyBytes32)
      expect(await nebula.subscriptionsNextElement(emptyBytes32)).to.eq(emptyBytes32)
      expect(await nebula.subscriptionsPrevElement(id1)).to.eq(emptyBytes32)
      expect(await nebula.subscriptionsNextElement(id1)).to.eq(id2)
      expect(await nebula.subscriptionsPrevElement(id2)).to.eq(id1)
      expect(await nebula.subscriptionsNextElement(id2)).to.eq(emptyBytes32)
    })

    it("updates subscriptions ids", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      expect(await nebula.subscriptionIds(0)).to.eq(id)
    })

    it("emits event", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      await expect(nebula.subscribe(subMockBytes.address, minConfirmations, reward))
        .to.emit(nebula, "NewSubscriber")
        .withArgs(id)
    })
  })

  describe("#sendValueToSubByte", () => {
    it("fails if there is no subscriber with subId", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await sendHashValue(ethers.utils.solidityKeccak256(["bytes"],[value]))

      await expect(nebula.connect(consul1).sendValueToSubByte(value, pulseId, id))
        .to.be.revertedWith("function call to a non-contract account")
    })

    it("fails if value was not approved by oracles", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      await expect(nebula.connect(consul1).sendValueToSubByte(value, pulseId, id))
        .to.be.revertedWith("value was not approved by oracles")
    })

    it("fails if a value was sent to subscriber this pulse", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["bytes"],[value]))

      await nebula.connect(consul1).sendValueToSubByte(value, 1, id)
      await expect(nebula.connect(consul1).sendValueToSubByte(value, pulseId, id))
        .to.be.revertedWith("sub sent")
    })

    it("updates isPulseSubSent", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      expect(await nebula.isPulseSubSent(pulseId, id)).to.eq(false)

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["bytes"],[value]))

      await nebula.connect(consul1).sendValueToSubByte(value, pulseId, id)

      expect(await nebula.isPulseSubSent(pulseId, id)).to.eq(true)
    })

    it("sends value to subscriber", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockBytes.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      expect(await subMockBytes.isSent()).to.eq(false)

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["bytes"],[value]))

      await nebula.connect(consul1).sendValueToSubByte(value, pulseId, id)

      expect(await subMockBytes.isSent()).to.eq(true)

    })
  })

  describe("#sendValueToSubInt", () => {
    it("fails if there is no subscriber with subId", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockInt.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await sendHashValue(ethers.utils.solidityKeccak256(["int64"],[value]))

      await expect(nebula.connect(consul1).sendValueToSubInt(value, pulseId, id))
        .to.be.revertedWith("function call to a non-contract account")
    })

    it("fails if value was not approved by oracles", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockInt.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      await expect(nebula.connect(consul1).sendValueToSubInt(value, pulseId, id))
        .to.be.revertedWith("value was not approved by oracles")
    })

    it("fails if a value was sent to subscriber this pulse", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockInt.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await nebula.subscribe(subMockInt.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["int64"],[value]))

      await nebula.connect(consul1).sendValueToSubInt(value, pulseId, id)
      await expect(nebula.connect(consul1).sendValueToSubInt(value, pulseId, id))
        .to.be.revertedWith("sub sent")
    })

    it("updates isPulseSubSent", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockInt.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      expect(await nebula.isPulseSubSent(pulseId, id)).to.eq(false)

      await nebula.subscribe(subMockInt.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["int64"],[value]))

      await nebula.connect(consul1).sendValueToSubInt(value, pulseId, id)

      expect(await nebula.isPulseSubSent(pulseId, id)).to.eq(true)
    })

    it("sends value to subscriber", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockInt.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      expect(await subMockInt.isSent()).to.eq(false)

      await nebula.subscribe(subMockInt.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["int64"],[value]))

      await nebula.connect(consul1).sendValueToSubInt(value, pulseId, id)

      expect(await subMockInt.isSent()).to.eq(true)

    })
  })

  describe("#sendValueToSubString", () => {
    it("fails if there is no subscriber with subId", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockString.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await sendHashValue(ethers.utils.solidityKeccak256(["string"],[value]))

      await expect(nebula.connect(consul1).sendValueToSubString(value, pulseId, id))
        .to.be.revertedWith("function call to a non-contract account")
    })

    it("fails if value was not approved by oracles", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockString.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await nebula.subscribe(subMockBytes.address, minConfirmations, reward)

      await expect(nebula.connect(consul1).sendValueToSubString(value, pulseId, id))
        .to.be.revertedWith("value was not approved by oracles")
    })

    it("fails if a value was sent to subscriber this pulse", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockString.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      await nebula.subscribe(subMockString.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["string"],[value]))

      await nebula.connect(consul1).sendValueToSubString(value, pulseId, id)
      await expect(nebula.connect(consul1).sendValueToSubString(value, pulseId, id))
        .to.be.revertedWith("sub sent")
    })

    it("updates isPulseSubSent", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockString.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      expect(await nebula.isPulseSubSent(pulseId, id)).to.eq(false)

      await nebula.subscribe(subMockString.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["string"],[value]))

      await nebula.connect(consul1).sendValueToSubString(value, pulseId, id)

      expect(await nebula.isPulseSubSent(pulseId, id)).to.eq(true)
    })

    it("sends value to subscriber", async () => {
      let minConfirmations = 1
      let reward = 0

      let sig = "0x3527715d"
      let pack1 = ethers.utils.solidityPack(["bytes4", "address", "address", "uint8"], [sig, wallet.address, subMockString.address, minConfirmations])
      let pack2 = ethers.utils.solidityPack(["bytes"], [pack1])
      let id = ethers.utils.solidityKeccak256(["bytes"], [pack2])

      let value = "0x01"
      let pulseId = 1

      expect(await subMockString.isSent()).to.eq(false)

      await nebula.subscribe(subMockString.address, minConfirmations, reward)

      await sendHashValue(ethers.utils.solidityKeccak256(["string"],[value]))

      await nebula.connect(consul1).sendValueToSubString(value, pulseId, id)

      expect(await subMockString.isSent()).to.eq(true)

    })
  })
})
