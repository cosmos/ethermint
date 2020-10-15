// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.6.7;

import "./CustodianUpgradeable.sol";

/**
 * @title Proxy
 * @dev Implements delegation of calls to other contracts, with proper
 * forwarding of return values and bubbling of failures.
 * It defines a fallback function that delegates all calls to the address
 * returned by the abstract _implementation() internal function.
 */
contract UpgradeableProxy is LockRequestable {
    address public proxyCustodian;
    address public proxyImpl;

    mapping (bytes32 => ProxyCustodianChangeRequest) public proxyCustodianChangeReqs;
    mapping (bytes32 => ProxyImplChangeRequest) public proxyImplChangeReqs;

    constructor(address _custodian, address _implementation) public {
        require(_custodian != address(0), "ERR_NO_ADDRESS");
        require(_isContract(_implementation), "ERR_NO_CONTRACT");
        proxyCustodian = _custodian;
        proxyImpl = _implementation;
    }

    function _isContract(address _addr) private view returns (bool){
        uint32 size;
        assembly {
            size := extcodesize(_addr)
        }
        return (size > 0);
    }

    modifier onlyProxyCustodian {
        require(msg.sender == proxyCustodian);
        _;
    }

    function requestProxyCustodianChange(address _proposedProxyCustodian) public returns (bytes32 lockId) {
        require(_proposedProxyCustodian != address(0));

        lockId = generateLockId();

        proxyCustodianChangeReqs[lockId] = ProxyCustodianChangeRequest({
            proposedNew: _proposedProxyCustodian
        });

        emit ProxyCustodianChangeRequested(lockId, msg.sender, _proposedProxyCustodian);
    }

    function confirmProxyCustodianChange(bytes32 _lockId) public onlyProxyCustodian {
        proxyCustodian = getProxyCustodianChangeReq(_lockId);

        delete proxyCustodianChangeReqs[_lockId];

        emit ProxyCustodianChangeConfirmed(_lockId, proxyCustodian);
    }

    function getProxyCustodianChangeReq(bytes32 _lockId) private view returns (address _proposedNew) {
        ProxyCustodianChangeRequest storage changeRequest = proxyCustodianChangeReqs[_lockId];

        // reject ‘null’ results from the map lookup
        // this can only be the case if an unknown `_lockId` is received
        require(changeRequest.proposedNew != address(0));

        return changeRequest.proposedNew;
    }

    /** @notice  Requests a change of the impl associated with this contract.
      *
      * @dev  Returns a unique lock id associated with the request.
      * Anyone can call this function, but confirming the request is authorized
      * by the impl.
      *
      * @param  _proposedProxyImpl  The address of the new impl.
      * @return  lockId  A unique identifier for this request.
      */
    function requestProxyImplChange(address _proposedProxyImpl) public returns (bytes32 lockId) {
        require(_isContract(_proposedProxyImpl), "ERR_NO_CONTRACT");

        lockId = generateLockId();

        proxyImplChangeReqs[lockId] = ProxyImplChangeRequest({
            proposedNew: _proposedProxyImpl
        });

        emit ProxyImplChangeRequested(lockId, msg.sender, _proposedProxyImpl);
    }

    /** @notice  Confirms a pending change of the implementation associated with this contract.
      *
      * @param  _lockId  The identifier of a pending change request.
      */
    function confirmProxyImplChange(bytes32 _lockId) public onlyProxyCustodian {
        proxyImpl = getProxyImplChangeReq(_lockId);

        delete proxyImplChangeReqs[_lockId];

        emit ProxyImplChangeConfirmed(_lockId, proxyImpl);
    }

    function getProxyImplChangeReq(bytes32 _lockId) private view returns (address _proposedNew) {
        ProxyImplChangeRequest storage changeRequest = proxyImplChangeReqs[_lockId];

        // reject ‘null’ results from the map lookup
        // this can only be the case if an unknown `_lockId` is received
        require(changeRequest.proposedNew != address(0));

        return changeRequest.proposedNew;
    }

    struct ProxyCustodianChangeRequest {
        address proposedNew;
    }

    event ProxyCustodianChangeRequested(
        bytes32 lockId,
        address msgSender,
        address proposedProxyCustodian
    );

    struct ProxyImplChangeRequest {
        address proposedNew;
    }

    event ProxyImplChangeRequested(
        bytes32 lockId,
        address msgSender,
        address proposedProxyImpl
    );

    event ProxyImplChangeConfirmed(bytes32 lockId, address newProxyImpl);
    event ProxyCustodianChangeConfirmed(bytes32 lockId, address newProxyCustodian);

    /**
     * @dev Delegates execution to an implementation contract.
     * This is a low level function that doesn't return to its internal call site.
     * It will return to the external caller whatever the implementation returns.
     * @param _implementation Address to delegate.
     */
    function _delegate(address _implementation) internal {
      assembly {
        // Copy msg.data. We take full control of memory in this inline assembly
        // block because it will not return to Solidity code. We overwrite the
        // Solidity scratch pad at memory position 0.
        calldatacopy(0, 0, calldatasize())

        // Call the implementation.
        // out and outsize are 0 because we don't know the size yet.
        let result := delegatecall(gas(), _implementation, 0, calldatasize(), 0, 0)

        // Copy the returned data.
        returndatacopy(0, 0, returndatasize())

        switch result
        // delegatecall returns 0 on error.
        case 0 { revert(0, returndatasize()) }
        default { return(0, returndatasize()) }
      }
    }

    fallback () payable external {
        _delegate(proxyImpl);
    }

    receive () payable external {
        _delegate(proxyImpl);
    }
}