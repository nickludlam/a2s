package a2s

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

// isMultiPacket checks if response uses multi-packet format.
func isMultiPacket(data []byte, showRawResponses bool) (bool, error) {
	// Some servers can send a truncated packet first; avoid panics on short reads.
	if len(data) < 4 {
		return true, ErrMultiPacket
	}

	header := binary.LittleEndian.Uint32(data[:4])

	switch header {
	case singlePacket:
		if len(data) < 5 {
			return false, ErrSinglePacket
		}
		return false, nil

	case multiPacket:
		if len(data) < 9 {
			return true, ErrMultiPacket
		}
		return true, nil

	default:
		err := errors.Join(ErrValidatorHeader, fmt.Errorf("0x%X", header))
		if showRawResponses {
			// Print hex dump to stderr highlighting the wrong header bytes (first 4 bytes)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprint(os.Stderr, hexDump(data, 0, 4))
		}
		return false, err
	}
}

// validateResponseType verifies response type matches the request type.
func validateResponseType(request, response Flag, data []byte, showRawResponses bool) error {
	var err error

	switch request {
	case InfoRequest:
		if response != infoResponseSource && response != infoResponseGoldSource {
			err = errors.Join(ErrValidatorInfo, fmt.Errorf("0x%X", response))
		}

	case PlayerRequest:
		if response != playerResponse {
			err = errors.Join(ErrValidatorPlayer, fmt.Errorf("0x%X", response))
		}

	case RulesRequest:
		if response != rulesResponse {
			err = errors.Join(ErrValidatorRules, fmt.Errorf("0x%X", response))
		}

	case PingRequest:
		if response != pingResponse {
			err = errors.Join(ErrValidatorPing, fmt.Errorf("0x%X", response))
		}

	case ChallengeRequest:
		if response != challengeResponse {
			err = errors.Join(ErrValidatorChallenge, fmt.Errorf("0x%X", response))
		}

	default:
		err = errors.Join(ErrValidatorRequest, fmt.Errorf("0x%X", request))
	}

	if err != nil && showRawResponses {
		// Print hex dump to stderr highlighting the wrong response byte (byte at offset 4)
		fmt.Fprintln(os.Stderr, err)
		// Highlight byte 4 (the response type byte)
		fmt.Fprint(os.Stderr, hexDump(data, 4, 5))
	}

	return err
}
