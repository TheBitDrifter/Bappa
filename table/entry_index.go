package table

import (
	numbers_util "github.com/TheBitDrifter/util/numbers"
)

var _ EntryIndex = &entryIndex{}

type entryIndex struct {
	currEntryID EntryID
	entries     []entry
	recyclable  []entry
}

func (ei *entryIndex) NewEntries(n, start int, tbl Table) ([]Entry, error) {
	if n <= 0 {
		return nil, BatchOperationError{Count: n}
	}
	amountRecyclable := min(len(ei.recyclable), n)
	newEntries := []Entry{}

	// First use recyclable entries
	for i := 0; i < amountRecyclable; i++ {
		e := entry{
			id:       ei.recyclable[i].ID(),
			recycled: ei.recyclable[i].Recycled() + 1,
			table:    tbl,
			index:    start + i,
		}
		globalIndex := int(e.ID()) - 1

		// Ensure the entries slice has enough capacity
		for globalIndex >= len(ei.entries) {
			ei.entries = append(ei.entries, entry{})
		}

		ei.entries[globalIndex] = e
		newEntries = append(newEntries, e)
	}

	// Remove used recyclable entries
	ei.recyclable = ei.recyclable[amountRecyclable:]

	// Create new entries for the remaining
	leftover := n - amountRecyclable
	for i := 0; i < leftover; i++ {
		ei.currEntryID++
		entry := entry{
			id:       ei.currEntryID,
			recycled: 0,
			table:    tbl,
			index:    start + i + amountRecyclable,
		}
		ei.entries = append(ei.entries, entry)
		newEntries = append(newEntries, entry)
	}

	return newEntries, nil
}

func (ei *entryIndex) Entry(i int) (Entry, error) {
	if i < 0 || i >= len(ei.entries) {
		return nil, AccessError{Index: i, UpperBound: len(ei.entries)}
	}
	entry := ei.entries[i]
	if entry.id == 0 {
		return nil, InvalidEntryAccessError{}
	}
	return entry, nil
}

func (ei *entryIndex) UpdateIndex(id EntryID, rowIndex int) error {
	entryIndex := int(id - 1)
	if entryIndex < 0 || entryIndex >= len(ei.entries) {
		return AccessError{Index: entryIndex, UpperBound: len(ei.entries) - 1}
	}
	e := ei.entries[entryIndex]
	newEntry := entry{
		id:       e.ID(),
		recycled: e.Recycled(),
		index:    rowIndex,
		table:    e.table,
	}
	ei.entries[entryIndex] = newEntry
	return nil
}

func (ei *entryIndex) RecycleEntries(ids ...EntryID) error {
	uniqueIDs := numbers_util.UniqueInts(entryIDs(ids).toInts())

	uniqCount := len(uniqueIDs)
	entriesCount := len(ei.entries)
	if uniqCount <= 0 || uniqCount >= entriesCount {
		return BatchDeletionError{Capacity: uniqCount, BatchOperationError: BatchOperationError{Count: uniqCount}}
	}

	for _, id := range ids {
		index := id - 1
		if ei.entries[index].ID() == 0 {
			continue
		}
		zeroEntry := entry{
			id:       0,
			recycled: ei.entries[index].Recycled(),
			index:    0,
		}
		recycledEntry := entry{
			id:       id,
			recycled: ei.entries[index].Recycled(),
			index:    0,
		}
		ei.recyclable = append(ei.recyclable, recycledEntry)
		ei.entries[index] = zeroEntry
	}
	return nil
}

func (ei *entryIndex) Reset() error {
	ei.entries = ei.entries[:0]
	ei.recyclable = ei.recyclable[:0]
	ei.currEntryID = 0
	return nil
}

func (ei *entryIndex) Entries() []Entry {
	entriesAsInterface := make([]Entry, len(ei.entries))
	for i, en := range ei.entries {
		entriesAsInterface[i] = en
	}
	return entriesAsInterface
}

func (ei *entryIndex) Recyclable() []Entry {
	recyclableEntriesAsInterface := make([]Entry, len(ei.recyclable))
	for i, en := range ei.recyclable {
		recyclableEntriesAsInterface[i] = en
	}
	return recyclableEntriesAsInterface
}

// Use with caution, primarily for deser
func (ei *entryIndex) ForceNewEntry(id int, recycled, tblIndex int, tbl Table) error {
	index := id - 1

	if index >= len(ei.entries) {
		amountNeeded := index + 1 - len(ei.entries)

		newEntries := make([]entry, amountNeeded)
		newRecycled := make([]entry, amountNeeded-1)

		ei.recyclable = append(ei.recyclable, newRecycled...)
		ei.entries = append(ei.entries, newEntries...)
		ei.currEntryID = EntryID(id)
	}

	ei.entries[index] = entry{
		id:       EntryID(id),
		table:    tbl,
		index:    tblIndex,
		recycled: recycled,
	}

	return nil
}
