// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.6.7;

import "./LockRequestable.sol";
import "./Upgradeable.sol";

/** @title  A contract to inherit upgradeable custodianship.
  *
  * @notice  A contract that provides re-usable code for upgradeable
  * custodianship. That custodian may be an account or another contract.
  *
  * @dev  This contract is intended to be inherited by any contract
  * requiring a custodian to control some aspect of its functionality.
  * This contract provides the mechanism for that custodianship to be
  * passed from one custodian to the next.
  */
contract CustodianUpgradeable is Upgradeable, LockRequestable {

    // TYPES
    /// @dev  The struct type for pending custodian changes.
    struct CustodianChangeRequest {
        address proposedNew;
    }

    // MEMBERS
    /// @dev  The address of the account or contract that acts as the custodian.
    address public custodian;

    /// @dev  The map of lock ids to pending custodian changes.
    mapping (bytes32 => CustodianChangeRequest) public custodianChangeReqs;

    function initialize(address _custodian) public initializer {
        custodian = _custodian;
    }

    // MODIFIERS
    modifier onlyCustodian {
        require(msg.sender == custodian);
        _;
    }

    /** @notice  Requests a change of the custodian associated with this contract.
      *
      * @dev  Returns a unique lock id associated with the request.
      * Anyone can call this function, but confirming the request is authorized
      * by the custodian.
      *
      * @param  _proposedCustodian  The address of the new custodian.
      * @return  lockId  A unique identifier for this request.
      */
    function requestCustodianChange(address _proposedCustodian) public returns (bytes32 lockId) {
        require(_proposedCustodian != address(0));

        lockId = generateLockId();

        custodianChangeReqs[lockId] = CustodianChangeRequest({
            proposedNew: _proposedCustodian
        });

        emit CustodianChangeRequested(lockId, msg.sender, _proposedCustodian);
    }

    /** @notice  Confirms a pending change of the custodian associated with this contract.
      *
      * @dev  When called by the current custodian with a lock id associated with a
      * pending custodian change, the `address custodian` member will be updated with the
      * requested address.
      *
      * @param  _lockId  The identifier of a pending change request.
      */
    function confirmCustodianChange(bytes32 _lockId) public onlyCustodian {
        custodian = getCustodianChangeReq(_lockId);

        delete custodianChangeReqs[_lockId];

        emit CustodianChangeConfirmed(_lockId, custodian);
    }

    // PRIVATE FUNCTIONS
    function getCustodianChangeReq(bytes32 _lockId) private view returns (address _proposedNew) {
        CustodianChangeRequest storage changeRequest = custodianChangeReqs[_lockId];

        // reject ‘null’ results from the map lookup
        // this can only be the case if an unknown `_lockId` is received
        require(changeRequest.proposedNew != address(0));

        return changeRequest.proposedNew;
    }

    /// @dev  Emitted by successful `requestCustodianChange` calls.
    event CustodianChangeRequested(
        bytes32 lockId,
        address msgSender,
        address proposedCustodian
    );

    /// @dev Emitted by successful `confirmCustodianChange` calls.
    event CustodianChangeConfirmed(bytes32 lockId, address newCustodian);
}
