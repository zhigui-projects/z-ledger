package ascc

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/aclmgmt"
	"github.com/hyperledger/fabric/core/aclmgmt/resources"
	"github.com/hyperledger/fabric/core/policy"
)

const (
	//invoke functions
	ArchiveByDate string = "archiveByDate"

	// keys
	ArchiveByTxDateKey string = "archive_by_tx_date"
	ChaincodeName      string = "ascc"
)

var logger = flogging.MustGetLogger("ascc")

type ArchiveSysCC struct {
	aclProvider aclmgmt.ACLProvider
	// PolicyChecker is the interface used to perform
	// access control
	policyChecker policy.PolicyChecker
}

// New returns an instance of ASCC.
// Typically this is called once per peer.
func New(aclProvider aclmgmt.ACLProvider, policyChecker policy.PolicyChecker) *ArchiveSysCC {
	return &ArchiveSysCC{
		aclProvider:   aclProvider,
		policyChecker: policyChecker,
	}
}

func (a *ArchiveSysCC) Name() string { return ChaincodeName }

func (a *ArchiveSysCC) Chaincode() shim.Chaincode { return a }

// Init initializes ASCC
func (a *ArchiveSysCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Init ASCC")
	return shim.Success(nil)
}

// Invoke is called with args[0] contains the query function name, args[1]
// contains the chain ID, which is temporary for now until it is part of stub.
// Each function requires additional parameters as described below:
// # ArchiveByDate: Archive the block files which contains tx date before the specified date in args[2]
func (a *ArchiveSysCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()

	if len(args) < 2 {
		logger.Errorf("ASCC - Incorrect number of arguments, %d", len(args))
		return shim.Error(fmt.Sprintf("Incorrect number of arguments, %d", len(args)))
	}

	function := string(args[0])
	cid := string(args[1])

	logger.Infof("ASCC - invoke function: %s on chain: %s", function, cid)

	// Handle ACL:
	// 1. get the signed proposal
	sp, err := stub.GetSignedProposal()
	if err != nil {
		logger.Errorf("ASCC - GetSignedProposal got error: %s", err)
		return shim.Error(fmt.Sprintf("Failed retrieving signed proposal on executing %s with error %s", function, err))
	}

	switch function {
	case ArchiveByDate:
		if function != ArchiveByDate && len(args) < 3 {
			logger.Errorf("ASCC - missing 3rd arg for: %s", function)
			return shim.Error(fmt.Sprintf("missing 3rd argument for %s", function))
		}
		// 2. check archive by date policy
		if err = a.aclProvider.CheckACL(resources.Ascc_ArchiveByDate, cid, sp); err != nil {
			logger.Errorf("ASCC - access denied for [%s]: %s", function, err)
			return shim.Error(fmt.Sprintf("access denied for [%s]: %s", function, err))
		}
		return archiveByDate(stub, args[2])
	default:
		logger.Errorf("ASCC - requested function %s not found in ASCC.", function)
		return shim.Error(fmt.Sprintf("Requested function %s not found in ASCC.", function))
	}
}

func archiveByDate(stub shim.ChaincodeStubInterface, date []byte) pb.Response {
	if err := stub.PutState(ArchiveByTxDateKey, date); err != nil {
		logger.Errorf("ASCC - archiveByDate put state got error: %s", err)
		return shim.Error(fmt.Sprintf("ASCC put state got error: %s", err))
	}
	return shim.Success(nil)
}
