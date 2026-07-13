package store

// ApplyUnreadMention idempotently adjusts a chat's unread-mention count from a
// per-message signal. Marking a not-yet-tracked message increments the count;
// clearing a tracked message decrements it (floored at 0). Returns true when the
// count changed. No-op (false) for unknown chats.
func (s *SQLiteStore) ApplyUnreadMention(chatID int64, msgID int, hasMention bool) bool {
	s.mu.Lock()
	c, ok := s.chats[chatID]
	if !ok {
		s.mu.Unlock()
		return false
	}
	tracked := s.unreadMentionMsgs[chatID]
	_, isTracked := tracked[msgID]

	changed := false
	switch {
	case hasMention && !isTracked:
		if tracked == nil {
			tracked = make(map[int]struct{})
			s.unreadMentionMsgs[chatID] = tracked
		}
		tracked[msgID] = struct{}{}
		c.UnreadMentionsCount++
		changed = true
	case !hasMention && isTracked:
		delete(tracked, msgID)
		if c.UnreadMentionsCount > 0 {
			c.UnreadMentionsCount--
		}
		changed = true
	}
	if changed {
		s.chats[chatID] = c
	}
	s.mu.Unlock()
	if changed {
		s.persistChat(c)
	}
	return changed
}

// SetChatMentionsRead clears a chat's unread-mention count and its tracked
// message set, then persists. No-op for unknown chats.
func (s *SQLiteStore) SetChatMentionsRead(chatID int64) {
	s.mu.Lock()
	c, ok := s.chats[chatID]
	if !ok {
		s.mu.Unlock()
		return
	}
	delete(s.unreadMentionMsgs, chatID)
	c.UnreadMentionsCount = 0
	s.chats[chatID] = c
	s.mu.Unlock()
	s.persistChat(c)
}
