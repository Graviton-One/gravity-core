import { ethers, waffle } from "hardhat"
import { Gravity } from "../typechain/Gravity"
import { expect } from "./shared/expect"
import { gravityFixture } from "./shared/fixtures"

describe("Gravity", () => {
  const [wallet, other, consul1, consul2, consul3, consul4, consul5, consul6] = waffle.provider.getWallets()

  let gravity: Gravity

  let loadFixture: ReturnType<typeof waffle.createFixtureLoader>

  before("create fixture loader", async () => {
    loadFixture = waffle.createFixtureLoader([wallet, other, consul1, consul2, consul3, consul4, consul5, consul6])
  })

  beforeEach("deploy test contracts", async () => {
    ;({ gravity } = await loadFixture(gravityFixture))
  })

  function packAddresses(addresses: string[]): string {
    var hash: string = "0x"
    for (var i in addresses) {
      hash = ethers.utils.solidityPack([ "bytes", "address" ], [ hash, addresses[i] ]);
    }
    return hash
  }

  function hashAddresses(addresses: string[], roundId: number): string {
    let hash = packAddresses(addresses)
    return ethers.utils.solidityKeccak256([ "bytes", "uint256" ], [ hash, roundId ])
  }

  it("constructor initializes variables", async () => {
    expect(await gravity.bftValue()).to.eq(3)
    expect(await gravity.rounds(0, 0)).to.eq(consul1.address)
    expect(await gravity.rounds(0, 1)).to.eq(consul2.address)
    expect(await gravity.rounds(0, 2)).to.eq(consul3.address)
  })

  it("starting state after deployment", async () => {
    expect(await gravity.lastRound()).to.eq(0)
  })

  describe("#getConsuls", () => {
    it("returns last round consuls", async () => {
      var consuls0 = await gravity.getConsuls()
      expect(consuls0.length).to.eq(3)
      expect(consuls0[0]).to.eq(consul1.address)
      expect(consuls0[1]).to.eq(consul2.address)
      expect(consuls0[2]).to.eq(consul3.address)
    })
  })

  describe("#getConsulsByRoundId", () => {
    it("returns last round consuls when roundId is 0", async () => {
      var consuls0 = await gravity.getConsulsByRoundId(0)
      expect(consuls0.length).to.eq(3)
      expect(consuls0[0]).to.eq(consul1.address)
      expect(consuls0[1]).to.eq(consul2.address)
      expect(consuls0[2]).to.eq(consul3.address)
    })
    it("returns consuls in a given round", async () => {
      // await gravity.updateConsuls(newConsuls, v, r, roundId)
      // expect(await gravity.getConsulsByRoundId(1)).to.eq([consul1.address, consul2.address, consul3.address])
    })
  })

  describe("#hashNewConsuls", () => {
    it("hashes one address", async () => {
      let roundId = 1
      let pack = ethers.utils.solidityPack([ "address" ], [ consul2.address ]);
      let hash = ethers.utils.solidityKeccak256([ "bytes", "uint256" ], [ pack, roundId ])

      let hashNewConsuls = await gravity.hashNewConsuls([consul2.address], 1)

      expect(hashNewConsuls).to.eq(hash)
    })

    it("hashes three addresses", async () => {
      let roundId = 1
      let pack1 = ethers.utils.solidityPack([ "address" ], [ consul1.address ]);
      let pack2 = ethers.utils.solidityPack([ "bytes", "address" ], [ pack1, consul2.address ]);
      let pack3 = ethers.utils.solidityPack([ "bytes", "address" ], [ pack2, consul3.address ]);
      let hash = ethers.utils.solidityKeccak256([ "bytes", "uint256" ], [ pack3, roundId ])

      let hashNewConsuls = await gravity.hashNewConsuls([consul1.address, consul2.address, consul3.address], 1)

      expect(hashNewConsuls).to.eq(hash)
    })

    it("hashes three addresses", async () => {
      let roundId = 1
      let addresses = [consul1.address, consul2.address, consul3.address]
      let pack = packAddresses(addresses)
      let hash = ethers.utils.solidityKeccak256([ "bytes", "uint256" ], [ pack, roundId ])

      let hashNewConsuls = await gravity.hashNewConsuls(addresses, 1)

      expect(hashNewConsuls).to.eq(hash)
    })

    it("hashes three addresses", async () => {
      let roundId = 1
      let addresses = [consul1.address, consul2.address, consul3.address]
      let hash = hashAddresses(addresses, roundId)

      let hashNewConsuls = await gravity.hashNewConsuls(addresses, 1)

      expect(hashNewConsuls).to.eq(hash)
    })
  })

  describe("#updateConsuls", () => {
    it("fails if round is less then last round", async () => {
      let roundId = 0
      let addresses = [consul4.address, consul5.address, consul6.address]
      let hash = hashAddresses(addresses, roundId)

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

      await expect(gravity.updateConsuls(addresses, vs, rs, ss, roundId))
        .to.be.revertedWith("round less last round");
    })

    it("fails if new oracles are not signed by at bft number oracles", async () => {
      let roundId = 1
      let addresses = [consul4.address, consul5.address, consul6.address]
      let hash = hashAddresses(addresses, roundId)

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

      await expect(gravity.updateConsuls(addresses, vs, rs, ss, roundId))
        .to.be.revertedWith("invalid bft count");
    })

    it("updates consuls", async () => {
      let roundId = 1
      let addresses = [consul4.address, consul5.address, consul6.address]
      let hash = hashAddresses(addresses, roundId)

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

      await gravity.updateConsuls(addresses, vs, rs, ss, roundId);

      expect(await gravity.rounds(1, 0)).to.eq(consul4.address)
      expect(await gravity.rounds(1, 1)).to.eq(consul5.address)
      expect(await gravity.rounds(1, 2)).to.eq(consul6.address)
    })
  })

})
