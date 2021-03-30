val hashValueScript =
    s"""{
        | // We expect msgHash to be in R4
        | val msgHash = OUTPUTS(0).R4[Coll[Byte]].get
        |
        | // We expect first option of signs to be in R6 [a, a, ..] TODO: after fix AOT in ergo this can be change to [(a, z), (a, z), ...]
        | val signs_a = OUTPUTS(0).R5[Coll[GroupElement]].get
        | // We expect first option of signs to be in R7 [z, z, ..]
        | val signs_z = OUTPUTS(0).R6[Coll[BigInt]].get
        |
        | // We expect pulseId to be in R6 and increase pulseId in out box
        | val checkPulse = OUTPUTS(0).R7[BigInt].get == SELF.R6[BigInt].get + 1
        |
        | val dataInput = CONTEXT.dataInputs(0)
        |
        | val check_NFT_tokens = {  allOf(Coll(
        |   // We expect one tokenNFT for hashValue contract to be in token(0)
        |   OUTPUTS(0).tokens(0)._1 == SELF.tokens(0)._1,
        |   // We expect one tokenNFT for oracle contract to be in token(0) of this box
        |   dataInput.tokens(0)._1 == oracleNebulaNFT // ‌Build Time
        | ))}
        | // get BftCoefficient from R4 of oracleContract Box
        | val bftValue = dataInput.R4[Int].get
        |
        | // Get oracles from R5 of oracleContract Box and convert to Coll[GroupElement]
        | val oracles: Coll[GroupElement] = dataInput.R5[Coll[Coll[Byte]]].get.map({ (oracle: Coll[Byte]) =>
        |     decodePoint(oracle)
        | })
        |
        | // Verify signs
        | val validateSign = {(v: ((Coll[Byte], GroupElement), (GroupElement, BigInt))) => {
        |    val e: Coll[Byte] = blake2b256(v._1._1) // weak Fiat-Shamir
        |    val eInt = byteArrayToBigInt(e) // challenge as big integer
        |    val g: GroupElement = groupGenerator
        |    val l = g.exp(v._2._2)
        |    val r = v._2._1.multiply(v._1._2.exp(eInt))
        |    if (l == r) 1 else 0
        | }}
        |
        | val count = validateSign(((msgHash, oracles(0)),(signs_a(0), signs_z(0))))
        |   + validateSign(((msgHash, oracles(1)),(signs_a(1), signs_z(1))))
        |   + validateSign(((msgHash, oracles(2)),(signs_a(2), signs_z(2))))
        |   + validateSign(((msgHash, oracles(3)),(signs_a(3), signs_z(3))))
        |   + validateSign(((msgHash, oracles(4)),(signs_a(4), signs_z(4))))
        |
        | // We Expect number of oracles that verified msgHash of in pulseId bigger than bftValue
        | val check_bftCoefficient = count >= bftValue
        |
        | sigmaProp (checkPulse && check_NFT_tokens && check_bftCoefficient)
        | }
    """.stripMargin


   val oracleScript =
   s"""{
      | // We get bftCoefficient from R4
      | val bftValue = SELF.R4[Int].get
      | val bftValueOut = OUTPUT(0).R4[Int].get
      |
      | // We get oracles from R5
      | val newSortedOracles = OUTPUTS(0).R5[Coll[Coll[Byte]]].get
      |
      | // We expect five oracles at R5 of OUTPUTS(0).
      | val check_oracles = newSortedOracles.size == 5
      |
      | // Verify signs
      | val validateSign = {(v: ((Coll[Byte], GroupElement), (GroupElement, BigInt))) => {
      |    val e: Coll[Byte] = blake2b256(v._1._1) // weak Fiat-Shamir
      |    val eInt = byteArrayToBigInt(e) // challenge as big integer
      |    val g: GroupElement = groupGenerator
      |    val l = g.exp(v._2._2)
      |    val r = v._2._1.multiply(v._1._2.exp(eInt))
      |    if (l == r) 1 else 0
      | }}
      |
      | // We expect first option of signs to be in R6 [a, a, ..] TODO: after fix AOT in ergo this can be change to [(a, z), (a, z), ...]
      | val signs_a = OUTPUTS(0).R6[Coll[GroupElement]].get
      | // We expect first option of signs to be in R7 [z, z, ..]
      | val signs_z = OUTPUTS(0).R7[Coll[BigInt]].get
      |
      | // should to be box of gravity contract
      | val dataInput = CONTEXT.dataInputs(0)
      |
      | val check_NFT_tokens = {  allOf(Coll(
      |   // We expect a NFT token for oracle contract to be in token 0
      |   OUTPUTS(0).tokens(0)._1 == SELF.tokens(0)._1
      |   // We expect in tokens of gravity contract there is NFT token of gravity
      |   dataInput.tokens(0)._1 == gravityNFT // ‌Build Time. TODO: Check gravityNFT state in gravity contract
      | ))}
      |
      | // We expect in R5 of gravity contract there are consuls
      | // TODO: Check register number in gravity contract
      | val consuls: Coll[GroupElement] = dataInput.R5[Coll[Coll[Byte]]].get.map({ (consul: Coll[Byte]) =>
      |     decodePoint(consul)
      | })
      |
      | // Concatenation all new oracles for create, newSortedOracles as a Coll[Byte] and verify signs.
      | val newSortedOracles1 = newSortedOracles.fold(fromBase64(""), { (baseOracle: Coll[Byte], newOracle: Coll[Byte]) => baseOracle ++ newOracle })
      | val count = validateSign(((newSortedOracles1, consuls(0)),(signs_a(0), signs_z(0))))
      |   + validateSign(((newSortedOracles1, consuls(1)),(signs_a(1), signs_z(1))))
      |   + validateSign(((newSortedOracles1, consuls(2)),(signs_a(2), signs_z(2))))
      |   + validateSign(((newSortedOracles1, consuls(3)),(signs_a(3), signs_z(3))))
      |   + validateSign(((newSortedOracles1, consuls(4)),(signs_a(4), signs_z(4))))
      |
      | // We Expect number of consuls that verified new oracles list bigger than bftValue
      | val check_bftCoefficient = { allOf(Coll(
      |   count >= bftValue,
      |   bftValueOut <= 5,
      |   bftValueOut > 0
      | ))}
      |
      | sigmaProp (check_NFT_tokens && check_bftCoefficient && check_oracles)
      |
      | }
   """.stripMargin
