pragma solidity ^0.8.0;

contract Voting{
    // mapping 来存储候选人的得票数
    mapping(string =>uint256) public candidateVotes;

    // 存储所有候选人的数组，用于重置功能
    string[] private candidates;

    // 投票事件，用于记录投票历史
    event Voted(address indexed voter,string candidate, uint256 newVoteCount);

    // 重置事件
    event VotesReset(address indexed resetBy);

    /**
     * @dev 投票给某个候选人
     * @param candidate 候选人名称
     */
    function vote(string memory candidate)public  {
        require(bytes(candidate).length>0,"Candidate name cannot be empty");

        // 如果是第一次给这个候选人投票，将其添加到候选人数组
        if(candidateVotes[candidate] == 0){
            candidates.push(candidate);
        }

        candidateVotes[candidate]++;
        emit Voted(msg.sender,candidate,candidateVotes[candidate]);
    }

    /**
     * @dev 获取某个候选人的得票数
     * @param candidate 候选人名称
     * @return 该候选人的得票数
     */

    function getVotes(string memory candidate)public view returns(uint256){
        return  candidateVotes[candidate];

    }
    /**
     * @dev 重置所有候选人的得票数
     */
    function resetVotes() public  {
        for(uint256 i=0;i<candidates.length;i++){
            candidateVotes[candidates[i]] = 0;
        }
        emit VotesReset(msg.sender);
    }
    /**
     * @dev 获取所有候选人列表（额外功能）
     * @return 包含所有候选人名称的数组
     */
    function getAllCandidates()public view returns(string[]  memory)  {
        return candidates;
    }

    /**
     * @dev 获取候选人总数（额外功能）
     * @return 候选人总数
     */
    function getCandidateCount() public view  returns(uint256) {
        return candidates.length;
    }

}