val gravityScript =
      s"""{
            |  val newConsuls = OUTPUTS(0).R5[Coll[Coll[Byte]]].get
            |  // make Coll[GroupElement] for sign validation from input's consuls witch are in [Coll[Coll[Byte]]] format
            |  val consuls: Coll[GroupElement] = SELF.R5[Coll[Coll[Byte]]].get.map({(consul: Coll[Byte]) => decodePoint(consul)})
            |
            |  // each sign made two part a (a groupelemet) and z(a bigint)
            |  val signs_a = OUTPUTS(0).R6[Coll[GroupElement]].get
            |  val signs_z = OUTPUTS(0).R7[Coll[BigInt]].get
            |  // making the message by concatenation of newConsoles
            |  val msg = newConsuls(0) ++ newConsuls(1) ++ newConsuls(2) ++ newConsuls(3) ++ newConsuls(4) })
            |
            | // Verify sign base on schnorr protocol
            |  val validateSign = {(v: ((Coll[Byte], GroupElement), (GroupElement, BigInt))) => {
            |     val e: Coll[Byte] = blake2b256(v._1._1) // weak Fiat-Shamir
            |     val eInt = byteArrayToBigInt(e) // challenge as big integer
            |     val g: GroupElement = groupGenerator
            |     val l = g.exp(v._2._2)
            |     val r = v._2._1.multiply(v._1._2.exp(eInt))
            |     if (l == r) 1 else 0
            |  }}
            |
            |  // validate each sign and consul
            |  val count = validateSign( ( (msg, consuls(0)), (signs_a(0), signs_z(0)) ) ) +
            |              validateSign( ( (msg, consuls(1)), (signs_a(1), signs_z(1)) ) ) +
            |              validateSign( ( (msg, consuls(2)), (signs_a(2), signs_z(2)) ) ) +
            |              validateSign( ( (msg, consuls(3)), (signs_a(3), signs_z(3)) ) ) +
            |              validateSign( ( (msg, consuls(4)), (signs_a(4), signs_z(4)) ) )
            |
            |  sigmaProp (
            |    allOf(Coll(
            |       // check output's bftvalue be valid
            |      OUTPUTS(0).R4[Int].get > 0 &&  OUTPUTS(0).R4[Int].get <= 5,
            |      OUTPUTS(0).propositionBytes == SELF.propositionBytes,
            |
            |      OUTPUTS(0).tokens(0)._1 == tokenId, // Build-time assignment, it's the NFT tocken
            |      OUTPUTS(0).tokens(0)._2 == 1,       // check NFT count
            |      OUTPUTS(0).value >= SELF.value,     // value of output should be bigger or equal to input's
            |
            |       // check count be bigger than input's bftvalue. to change the consuls,
            |       // it's important to sign at least equal to input's bftvalue
            |      count >= SELF.R4[Int].get
            |
            |  )))
            |}""".stripMargin
