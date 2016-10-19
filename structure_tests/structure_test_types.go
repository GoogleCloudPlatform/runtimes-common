package main

type StructureTest interface {
	GetCommandTests() []CommandTest
	GetFileExistenceTests() []FileExistenceTest
	GetFileContentTests() []FileContentTest
}

type StructureTestv0 struct {
	CommandTests       []CommandTestv0
	FileExistenceTests []FileExistenceTestv0
	FileContentTests   []FileContentTestv0
}

func (t StructureTestv0) GetCommandTests() []CommandTest {
	var interfaceSlice []CommandTest = make([]CommandTest, len(t.CommandTests))
	for i, d := range t.CommandTests {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

func (t StructureTestv0) GetFileExistenceTests() []FileExistenceTest {
	var interfaceSlice []FileExistenceTest = make([]FileExistenceTest, len(t.FileExistenceTests))
	for i, d := range t.FileExistenceTests {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

func (t StructureTestv0) GetFileContentTests() []FileContentTest {
	var interfaceSlice []FileContentTest = make([]FileContentTest, len(t.FileContentTests))
	for i, d := range t.FileContentTests {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

type StructureTestv1 struct {
	CommandTests       []CommandTestv1
	FileExistenceTests []FileExistenceTestv0
	FileContentTests   []FileContentTestv0
}

func (t StructureTestv1) GetCommandTests() []CommandTest {
	var interfaceSlice []CommandTest = make([]CommandTest, len(t.CommandTests))
	for i, d := range t.CommandTests {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

func (t StructureTestv1) GetFileExistenceTests() []FileExistenceTest {
	var interfaceSlice []FileExistenceTest = make([]FileExistenceTest, len(t.FileExistenceTests))
	for i, d := range t.FileExistenceTests {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

func (t StructureTestv1) GetFileContentTests() []FileContentTest {
	var interfaceSlice []FileContentTest = make([]FileContentTest, len(t.FileContentTests))
	for i, d := range t.FileContentTests {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}
