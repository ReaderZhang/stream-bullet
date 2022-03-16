package logic

//删除用户
func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var empty_Slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			empty_Slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return empty_Slice
}
