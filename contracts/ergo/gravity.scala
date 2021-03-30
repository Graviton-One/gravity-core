val gravityScript =
      s"""{
         |  val bftValue = SELF.R4[Int].get
         |  val newConsuls = OUTPUTS(0).R5[Coll[Coll[Byte]]].get
         |
         |  val consuls: Coll[GroupElement] = SELF.R5[Coll[Coll[Byte]]].get.map({(consul: Coll[Byte]) => decodePoint(consul)})
         |  val signs_a = OUTPUTS(0).R6[Coll[GroupElement]].get
         |  val signs_z = OUTPUTS(0).R7[Coll[BigInt]].get
         |
         |  val msg = newConsuls.fold(fromBase64(""), { (baseConsul: Coll[Byte], newConsul: Coll[Byte]) => baseConsul ++ newConsul })
         |
         | // Verify sign
         |  val validateSign = {(v: ((Coll[Byte], GroupElement), (GroupElement, BigInt))) => {
         |     val e: Coll[Byte] = blake2b256(v._1._1) // weak Fiat-Shamir
         |     val eInt = byteArrayToBigInt(e) // challenge as big integer
         |     val g: GroupElement = groupGenerator
         |     val l = g.exp(v._2._2)
         |     val r = v._2._1.multiply(v._1._2.exp(eInt))
         |     if (l == r) 1 else 0
         |  }}
         |
         |  val count = validateSign( ( (msg, consuls(0)), (signs_a(0), signs_z(0)) ) ) +
         |              validateSign( ( (msg, consuls(1)), (signs_a(1), signs_z(1)) ) )
         |
         |  sigmaProp (
         |    allOf(Coll(
         |      bftValue > 0,
         |      OUTPUTS(0).propositionBytes == SELF.propositionBytes,
         |      OUTPUTS(0).tokens(0)._1 == tokenId, // Build-time assignment
         |      OUTPUTS(0).tokens(0)._2 == 1,
         |      OUTPUTS(0).value >= SELF.value
         |      count >= bftValue
         |
         |  )))
         |}""".stripMargin
