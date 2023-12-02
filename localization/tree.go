package localization

type TrieNode struct {
	children map[rune]*TrieNode
	endOfKey bool
}

type Trie struct {
	root *TrieNode
}

func newTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

func (t *Trie) insert(key string) {
	node := t.root
	for _, char := range key {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.endOfKey = true
}

func (t *Trie) searchPrefix(prefix string) []string {
	node := t.root
	for _, char := range prefix {
		if _, ok := node.children[char]; !ok {
			return nil
		}
		node = node.children[char]
	}
	var results []string
	t.collectAllKeys(node, prefix, &results)
	return results
}

func (t *Trie) collectAllKeys(node *TrieNode, prefix string, results *[]string) {
	if node.endOfKey {
		*results = append(*results, prefix)
	}
	for char, childNode := range node.children {
		t.collectAllKeys(childNode, prefix+string(char), results)
	}
}
