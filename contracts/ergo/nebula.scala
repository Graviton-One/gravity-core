val signalScript: String =
    s"""{
       | sigmaProp(allOf(Coll(
       |  // There must be a msgHash in the R4 of the signal box
       |  SELF.R4[Coll[Byte]].isDefined,  // TODO: In the future, we have to check the msgHash for the USER-SC.
       |
       |  // Id of first token in signal box must be equal to tokenRepoId with value 1
       |  SELF.tokens(0)._1 == tokenRepoId,
       |  SELF.tokens(0)._2 == 1,
       |
       |  // Contract of second INPUT must be equal to tokenRepoContractHash
       |  blake2b256(INPUTS(1).propositionBytes) == tokenRepoContractHash,
       |  // Id of first token in token repo box must be equal to tokenRepoId
       |  INPUTS(1).tokens(0)._1 == tokenRepoId,
       |
       |  // Contract of first OUTPUT must be equal to tokenRepoContractHash
       |  blake2b256(OUTPUTS(0).propositionBytes) == tokenRepoContractHash
       | )))
       |}""".stripMargin

val tokenRepoScript: String =
    s"""{
       | val checkPulse = {allOf(Coll(
       |  // Contract of new tokenRepo box must be equal to contract of tokenRepo box in input
       |  SELF.propositionBytes == OUTPUTS(1).propositionBytes,
       |  // Id of first token in tokenRepo box must be equal to tokenRepoId
       |  SELF.tokens(0)._1 == tokenRepoId,
       |  // The transaction in which the tokenRepo box is located as the input box must contain the first input box containing the pulseNebulaNFT token
       |  INPUTS(0).tokens(0)._1 == pulseNebulaNFT,
       |  // OUTPUTS(1) is box of tokenRepo, OUTPUTS(2) is box of signal
       |  // In scenario add_pulse, a token is transferred from the tokenRepo to the signal box, also the minValue value must be sent to the signal box.
       |  OUTPUTS(1).tokens(0)._1 == tokenRepoId,
       |  OUTPUTS(1).tokens(0)._2 == SELF.tokens(0)._2 - 1,
       |  OUTPUTS(1).value == SELF.value - minValue,
       |  OUTPUTS(2).tokens(0)._1 == tokenRepoId,
       |  OUTPUTS(2).tokens(0)._2 == 1,
       |  OUTPUTS(2).value == minValue
       | ))}
       | // In scenario spend signal box in USER-SC, the token in the signal  box and its Erg must be returned to the tokenRepo.
       | val checkSignal = {allOf(coll(
       |  OUTPUTS(0).value == SELF.value + minValue,
       |  OUTPUTS(0).tokens(0)._1 == tokenRepoId,
       |  OUTPUTS(0).tokens(0)._2 == INPUTS(1).tokens(0)._2 + 1,
       |  OUTPUTS(0).propositionBytes == SELF.propositionBytes
       | ))}
       | sigmaProp(checkPulse || checkSignal)
       |}""".stripMargin

val pulseScript: String =
    s"""{
       | // We expect msgHash to be in R4
       | val msgHash = OUTPUTS(0).R4[Coll[Byte]].get
       |
       | // We expect first option of signs to be in R6 [a, a, ..] TODO: after fix AOT in ergo this can be change to [(a, z), (a, z), ...]
       | val signs_a = OUTPUTS(0).R5[Coll[GroupElement]].get
       | // We expect second option of signs to be in R7 [z, z, ..]
       | val signs_z = OUTPUTS(0).R6[Coll[BigInt]].get
       |
       | // should to be box of oracle contract
       | val dataInput = CONTEXT.dataInputs(0)
       |
       | // Verify signs
       | val validateSign: Int = {(v: ((Coll[Byte], GroupElement), (GroupElement, BigInt))) => {
       |    val e: Coll[Byte] = blake2b256(v._1._1) // weak Fiat-Shamir
       |    val eInt = byteArrayToBigInt(e) // challenge as big integer
       |    val g: GroupElement = groupGenerator
       |    val l = g.exp(v._2._2)
       |    val r = v._2._1.multiply(v._1._2.exp(eInt))
       |    if (l == r) 1 else 0
       | }}
       |
       | // We Expect number of oracles that verified msgHash of in pulseId bigger than bftValue
       | val check_bftCoefficient = {
       |   // We expect one tokenNFT for oracle contract to be in token(0) of this box
       |   if (dataInput.tokens(0)._1 == oracleNebulaNFT) {
       |     // get BftCoefficient from R4 of oracleContract Box
       |     val bftValue = dataInput.R4[Int].get
       |     // Get oracles from R5 of oracleContract Box and convert to Coll[GroupElement]
       |     val oracles: Coll[GroupElement] = dataInput.R5[Coll[Coll[Byte]]].get.map({ (oracle: Coll[Byte]) =>
       |         decodePoint(oracle)
       |     })
       |     val count : Int= validateSign(((msgHash, oracles(0)),(signs_a(0), signs_z(0)))) + validateSign(((msgHash, oracles(1)),(signs_a(1), signs_z(1)))) + validateSign(((msgHash, oracles(2)),(signs_a(2), signs_z(2)))) + validateSign(((msgHash, oracles(3)),(signs_a(3), signs_z(3)))) + validateSign(((msgHash, oracles(4)),(signs_a(4), signs_z(4))))
       |     count >= bftValue
       |   }
       |  else false
       | }
       |
       | val checkOUTPUTS = {
       |   if(SELF.tokens(0)._1 == pulseNebulaNFT) {
       |    allOf(Coll(
       |      // We expect one tokenNFT for pulse contract to be in token(0)
       |      OUTPUTS(0).tokens(0)._1 == pulseNebulaNFT,
       |      // Value of new pulse box must be greater than equal to value of pulse box in input
       |      OUTPUTS(0).value >= SELF.value,
       |      // Contract of new pulse box must be equal to contract of pulse box in input
       |      OUTPUTS(0).propositionBytes == SELF.propositionBytes,
       |      // We expect pulseId to be in R7 and increase pulseId in out box
       |      OUTPUTS(0).R7[BigInt].get == SELF.R7[BigInt].get + 1
       |
       |      // Contract of second INPUT/OUTPUT must be equal to tokenRepoContractHash
       |      blake2b256(INPUTS(1).propositionBytes) == tokenRepoContractHash,
       |      blake2b256(OUTPUTS(1).propositionBytes) == tokenRepoContractHash,
       |
       |      // Contract of third OUTPUT must be equal to signalContractHash
       |      blake2b256(OUTPUTS(2).propositionBytes) == signalContractHash,
       |      // There must be a msgHash in the R4 of the signal box
       |      OUTPUTS(2).R4[Coll[Byte]].get == msgHash
       |    ))
       |   }
       |   else false
       | }
       |
       | sigmaProp ( check_bftCoefficient && checkOUTPUTS )
       |
       | }
    """.stripMargin


val oracleScript: String =
    s"""{
       | // We get oracles from R5
       | val newSortedOracles = OUTPUTS(0).R5[Coll[Coll[Byte]]].get
       |
       | // We expect first option of signs to be in R6 [a, a, ..] TODO: after fix AOT in ergo this can be change to [(a, z), (a, z), ...]
       | val signs_a = OUTPUTS(0).R6[Coll[GroupElement]].get
       | // We expect first option of signs to be in R7 [z, z, ..]
       | val signs_z = OUTPUTS(0).R7[Coll[BigInt]].get
       |
       | // should to be box of gravity contract
       | val dataInput = CONTEXT.dataInputs(0)
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
       | val check_bftCoefficient = {
       |   // We expect in tokens of gravity contract there is NFT token of gravity also five oracles at R5 of OUTPUTS(0)
       |   if (dataInput.tokens(0)._1 == gravityNFT && newSortedOracles.size == 5) {
       |     // We get bftCoefficient from R4
       |     val bftValueIn = SELF.R4[Int].get
       |     val bftValueOut = OUTPUTS(0).R4[Int].get
       |     // We expect in R5 of gravity contract there are consuls
       |     val consuls: Coll[GroupElement] = dataInput.R5[Coll[Coll[Byte]]].get.map({ (consul: Coll[Byte]) =>
       |       decodePoint(consul)
       |     })
       |     // Concatenation all new oracles for create, newSortedOracles as a Coll[Byte] and verify signs.
       |     val newSortedOracles1 = newSortedOracles(0) ++ newSortedOracles(1) ++ newSortedOracles(2) ++ newSortedOracles(3) ++ newSortedOracles(4)
       |     val count = validateSign(((newSortedOracles1, consuls(0)),(signs_a(0), signs_z(0)))) + validateSign(((newSortedOracles1, consuls(1)),(signs_a(1), signs_z(1)))) + validateSign(((newSortedOracles1, consuls(2)),(signs_a(2), signs_z(2)))) + validateSign(((newSortedOracles1, consuls(3)),(signs_a(3), signs_z(3)))) + validateSign(((newSortedOracles1, consuls(4)),(signs_a(4), signs_z(4))))
       |     // We Expect the numbers of consuls that verified the new oracles list, to be more than three. TODO: in the future, with a change in the contract, this parameter can be dynamic.
       |     bftValueIn == bftValueOut && count >= bftValue
       |   }
       |  else false
       | }
       |
       | val checkOUTPUT = {
       |   if(SELF.tokens(0)._1 == oracleNebulaNFT) {
       |    allOf(Coll(
       |      // We expect a NFT token for oracle contract to be in tokens(0)
       |      OUTPUTS(0).tokens(0)._1 == oracleNebulaNFT,
       |
       |      // Value of new oracle box must be greater than equal to value of oracle box in input
       |      OUTPUTS(0).value >= SELF.value,
       |      // Contract of new oracle box must be equal to contract of oracle box in input
       |      OUTPUTS(0).propositionBytes == SELF.propositionBytes
       |    ))
       |   }
       |   else false
       | }
       |
       | sigmaProp ( checkOUTPUT && check_bftCoefficient )
       |
       | }
    """.stripMargin