// Wrapped in loadVueIfNeeded so the app works even if the layout
// did not load Vue (custom layout).
loadVueIfNeeded((err) => {
    if (err) { console.error('Vue load failed:', err); return; }
    const { createApp } = Vue;

    createApp({
    template: document.getElementById('app-group-members-template').innerHTML,
    data() {
        return {
            loading: true,
            addingMember: false,
            groupId: '',
            returnUrl: '',
            members: [],
            showAddModal: false,
            searchQuery: '',
            searchResults: [],
            searching: false,
            searchTimer: null
        };
    },
    mounted() {
        this.groupId = GROUP_ID_PLACEHOLDER;
        this.returnUrl = RETURN_URL_PLACEHOLDER;
        this.loadMembers();
    },
    methods: {
        statusClass(status) {
            switch (status) {
                case 'active': return 'bg-success';
                case 'inactive': return 'bg-danger';
                case 'unverified': return 'bg-info';
                case 'deleted': return 'bg-secondary';
                default: return 'bg-light text-dark';
            }
        },
        isMember(userId) {
            return this.members.some(m => m.user_id === userId);
        },
        async loadMembers() {
            this.loading = true;
            try {
                const response = await fetch(urlGetMembers, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: new URLSearchParams({
                        action: 'members-fetch-ajax',
                        group_id: this.groupId
                    })
                });
                const result = await response.json();
                if (result.status === 'success') {
                    this.members = result.data.members || [];
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to load members', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error loading members:', err);
                Notiflix.Notify.failure('Failed to load members', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                this.loading = false;
            }
        },
        openAddModal() {
            this.showAddModal = true;
            this.searchQuery = '';
            this.searchUsers();
            this.$nextTick(() => {
                if (this.$refs.searchInput) {
                    this.$refs.searchInput.focus();
                }
            });
        },
        closeAddModal() {
            this.showAddModal = false;
        },
        onSearchInput() {
            if (this.searchTimer) {
                clearTimeout(this.searchTimer);
            }
            this.searchTimer = setTimeout(() => this.searchUsers(), 300);
        },
        async searchUsers() {
            this.searching = true;
            try {
                const response = await fetch(urlUsersSearch, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: new URLSearchParams({
                        action: 'users-search-ajax',
                        group_id: this.groupId,
                        q: this.searchQuery.trim()
                    })
                });
                const result = await response.json();
                if (result.status === 'success') {
                    this.searchResults = result.data.users || [];
                }
            } catch (err) {
                console.error('Error searching users:', err);
            } finally {
                this.searching = false;
            }
        },
        async addMember(user) {
            this.addingMember = true;
            try {
                const response = await fetch(urlMemberAdd, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        group_id: this.groupId,
                        user_ref: user.user_id,
                        action: 'member-add-ajax'
                    })
                });
                const result = await response.json();
                if (result.status === 'success') {
                    Notiflix.Notify.success('Member added', { position: 'right-top', timeout: 3000 });
                    await this.loadMembers();
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to add member', { position: 'right-top', timeout: 3000 });
                }
            } catch (err) {
                console.error('Error adding member:', err);
                Notiflix.Notify.failure('Failed to add member', { position: 'right-top', timeout: 3000 });
            } finally {
                this.addingMember = false;
            }
        },
        async removeMember(member) {
            try {
                const response = await fetch(urlMemberRemove, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        group_id: this.groupId,
                        user_id: member.user_id,
                        action: 'member-remove-ajax'
                    })
                });
                const result = await response.json();
                if (result.status === 'success') {
                    Notiflix.Notify.success('Member removed', { position: 'right-top', timeout: 3000 });
                    await this.loadMembers();
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to remove member', { position: 'right-top', timeout: 3000 });
                }
            } catch (err) {
                console.error('Error removing member:', err);
                Notiflix.Notify.failure('Failed to remove member', { position: 'right-top', timeout: 3000 });
            }
        }
    }
}).mount('#app-group-members');
});
