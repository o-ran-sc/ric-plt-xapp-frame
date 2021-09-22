//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package writer

import (
	"fmt"
	rnibcommon "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	rnibentities "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/proto"
)

type rNibWriterInstance struct {
	sdl        rnibcommon.ISdlInstance //Deprecated: Will be removed in a future release and replaced by sdlStorage
	sdlStorage rnibcommon.ISdlSyncStorage
	ns         string
}

/*
RNibWriter interface allows saving data to the redis BD
*/
type RNibWriter interface {
	SaveNodeb(nbIdentity *rnibentities.NbIdentity, nb *rnibentities.NodebInfo) error
}

/*
GetNewRNibWriter returns reference to RNibWriter
*/
func GetNewRNibWriter(sdlStorage rnibcommon.ISdlSyncStorage) RNibWriter {
	return &rNibWriterInstance{
		sdl:        nil,
		sdlStorage: sdlStorage,
		ns:         rnibcommon.GetRNibNamespace(),
	}
}

//GetRNibWriter returns reference to RNibWriter
//Deprecated: Will be removed in a future release, please use GetNewRNibWriter instead.
func GetRNibWriter(sdl rnibcommon.ISdlInstance) RNibWriter {
	return &rNibWriterInstance{
		sdl:        sdl,
		sdlStorage: nil,
		ns:         "",
	}
}

/*
SaveNodeb saves nodeB entity data in the redis DB according to the specified data model
*/
func (w *rNibWriterInstance) SaveNodeb(nbIdentity *rnibentities.NbIdentity, entity *rnibentities.NodebInfo) error {
	isNotEmptyIdentity := isNotEmpty(nbIdentity)

	if isNotEmptyIdentity && entity.GetNodeType() == rnibentities.Node_UNKNOWN {
		return rnibcommon.NewValidationError(fmt.Sprintf("#rNibWriter.saveNodeB - Unknown responding node type, entity: %v", entity))
	}
	data, err := proto.Marshal(entity)
	if err != nil {
		return rnibcommon.NewInternalError(err)
	}
	var pairs []interface{}
	key, rNibErr := rnibcommon.ValidateAndBuildNodeBNameKey(nbIdentity.InventoryName)
	if rNibErr != nil {
		return rNibErr
	}
	pairs = append(pairs, key, data)

	if isNotEmptyIdentity {
		key, rNibErr = rnibcommon.ValidateAndBuildNodeBIdKey(entity.GetNodeType().String(), nbIdentity.GlobalNbId.GetPlmnId(), nbIdentity.GlobalNbId.GetNbId())
		if rNibErr != nil {
			return rNibErr
		}
		pairs = append(pairs, key, data)
	}

	if entity.GetEnb() != nil {
		pairs, rNibErr = appendEnbCells(nbIdentity.InventoryName, entity.GetEnb().GetServedCells(), pairs)
		if rNibErr != nil {
			return rNibErr
		}
	}
	if entity.GetGnb() != nil {
		pairs, rNibErr = appendGnbCells(nbIdentity.InventoryName, entity.GetGnb().GetServedNrCells(), pairs)
		if rNibErr != nil {
			return rNibErr
		}
	}
	if w.sdlStorage != nil {
		err = w.sdlStorage.Set(w.ns, pairs)
	} else {
		err = w.sdl.Set(pairs)
	}
	if err != nil {
		return rnibcommon.NewInternalError(err)
	}

	ranNameIdentity := &rnibentities.NbIdentity{InventoryName: nbIdentity.InventoryName}

	if isNotEmptyIdentity {
		nbIdData, err := proto.Marshal(ranNameIdentity)
		if err != nil {
			return rnibcommon.NewInternalError(err)
		}
		if w.sdlStorage != nil {
			err = w.sdlStorage.RemoveMember(w.ns, rnibentities.Node_UNKNOWN.String(), nbIdData)
		} else {
			err = w.sdl.RemoveMember(rnibentities.Node_UNKNOWN.String(), nbIdData)
		}
		if err != nil {
			return rnibcommon.NewInternalError(err)
		}
	} else {
		nbIdentity = ranNameIdentity
	}

	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		return rnibcommon.NewInternalError(err)
	}
	if w.sdlStorage != nil {
		err = w.sdlStorage.AddMember(w.ns, entity.GetNodeType().String(), nbIdData)
	} else {
		err = w.sdl.AddMember(entity.GetNodeType().String(), nbIdData)
	}
	if err != nil {
		return rnibcommon.NewInternalError(err)
	}
	return nil
}

/*
Close closes writer's pool
*/
func CloseWriter() {
}

func appendEnbCells(inventoryName string, cells []*rnibentities.ServedCellInfo, pairs []interface{}) ([]interface{}, error) {
	for _, cell := range cells {
		cellEntity := rnibentities.Cell{Type: rnibentities.Cell_LTE_CELL, Cell: &rnibentities.Cell_ServedCellInfo{ServedCellInfo: cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, rnibcommon.NewInternalError(err)
		}
		key, rNibErr := rnibcommon.ValidateAndBuildCellIdKey(cell.GetCellId())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = rnibcommon.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetPci())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}

func appendGnbCells(inventoryName string, cells []*rnibentities.ServedNRCell, pairs []interface{}) ([]interface{}, error) {
	for _, cell := range cells {
		cellEntity := rnibentities.Cell{Type: rnibentities.Cell_NR_CELL, Cell: &rnibentities.Cell_ServedNrCell{ServedNrCell: cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, rnibcommon.NewInternalError(err)
		}
		key, rNibErr := rnibcommon.ValidateAndBuildNrCellIdKey(cell.GetServedNrCellInformation().GetCellId())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = rnibcommon.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetServedNrCellInformation().GetNrPci())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}

func isNotEmpty(nbIdentity *rnibentities.NbIdentity) bool {
	return nbIdentity.GlobalNbId != nil && nbIdentity.GlobalNbId.PlmnId != "" && nbIdentity.GlobalNbId.NbId != ""
}
