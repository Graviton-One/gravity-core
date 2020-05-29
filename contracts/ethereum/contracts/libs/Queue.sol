pragma solidity >=0.5.16 <0.7.0;

library QueueLib {
    struct Queue {
        bytes32 first;
        bytes32 last;
        mapping(bytes32=>bytes32) nextElement;
        mapping(bytes32=>bytes32) prevElement;
    }

    function drop(Queue storage queue, bytes32 rqHash) public {
        if (queue.first == rqHash && queue.last == rqHash) {
            queue.first = 0x000;
            queue.last = 0x000;
        } else if (queue.first == rqHash) {
            queue.first = queue.nextElement[rqHash];
        } else if (queue.last == rqHash) {
            queue.last = queue.prevElement[rqHash];
        }
    }

    function next(Queue storage queue, bytes32 startRqHash) public view returns(bytes32) {
        if (startRqHash == 0x000)
            return queue.first;
        else {
            return queue.nextElement[startRqHash];
        }
    }


    function push(Queue storage queue, bytes32 elementHash) public {
        if (queue.first == 0x000) {
            queue.first = elementHash;
            queue.last = elementHash;
        } else {
            queue.nextElement[queue.last] = elementHash;
            queue.prevElement[elementHash] = queue.last;
            queue.last = elementHash;
        }
    }

}