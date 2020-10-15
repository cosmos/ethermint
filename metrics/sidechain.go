package metrics

func SidechainSubmitBatchSize(size int, tags ...Tags) {
	clientMux.RLock()
	defer clientMux.RUnlock()

	if client == nil {
		return
	}

	tagSpec := joinTags(tags...)

	client.Count("sidechain.submit_batch.size"+tagSpec, size)
}
