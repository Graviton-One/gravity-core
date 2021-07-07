import { ethers, waffle } from "hardhat"
import { Gravity } from "../typechain/Gravity"
import { expect } from "./shared/expect"
import { gravityFixture } from "./shared/fixtures"

describe("Gravity", () => {
  const [wallet, other, consul1, consul2, consul3] = waffle.provider.getWallets()

  let gravity: Gravity

  let loadFixture: ReturnType<typeof waffle.createFixtureLoader>

  before("create fixture loader", async () => {
    loadFixture = waffle.createFixtureLoader([wallet, other, consul1, consul2, consul3])
  })

  beforeEach("deploy test contracts", async () => {
    ;({ gravity } = await loadFixture(gravityFixture))
  })

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
    it("", async () => {
    })
  })

  describe("#getConsulsByRoundId", () => {
    it("", async () => {
    })
  })

  describe("#updateConsuls", () => {
    it("", async () => {
    })
  })

  describe("#hashNewConsuls", () => {
    it("", async () => {
    })
  })
})
