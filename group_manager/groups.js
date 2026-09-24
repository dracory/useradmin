const GroupsApp = {
    data() {
        return {
            loading: true,
            showCreateModal: false,
            creating: false,
            groups: [],
            urlGroupUpdate: urlGroupUpdate,
            newGroup: {
                name: '',
                handle: ''
            }
        };
    },
    mounted() {
        this.loadGroups();
    },
    methods: {
        async loadGroups() {
            this.loading = true;
            try {
                const response = await fetch(urlGroupsLoad, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({})
                });
                const data = await response.json();

                if (data.status === 'success') {
                    this.groups = data.data?.groups || [];
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to load groups'
                    });
                }
            } catch (error) {
                console.error('Error loading groups:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'Failed to load groups'
                });
            } finally {
                this.loading = false;
            }
        },
        statusClass(status) {
            switch (status) {
                case 'active': return 'bg-success';
                case 'inactive': return 'bg-danger';
                case 'deleted': return 'bg-secondary';
                default: return 'bg-light text-dark';
            }
        },
        async deleteGroup(group) {
            const result = await Swal.fire({
                icon: 'warning',
                title: 'Delete Group?',
                text: `Are you sure you want to delete ${group.name}?`,
                showCancelButton: true,
                confirmButtonText: 'Delete',
                cancelButtonText: 'Cancel',
                confirmButtonColor: '#dc3545'
            });

            if (!result.isConfirmed) return;

            try {
                const response = await fetch(urlGroupDelete, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ group_id: group.id })
                });

                const data = await response.json();

                if (data.status === 'success') {
                    Swal.fire({
                        icon: 'success',
                        title: 'Deleted',
                        text: 'Group deleted successfully',
                        timer: 1500,
                        showConfirmButton: false
                    });
                    this.loadGroups();
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to delete group'
                    });
                }
            } catch (error) {
                console.error('Error deleting group:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to delete group'
                });
            }
        },
        async createGroup() {
            if (!this.newGroup.name.trim()) {
                Swal.fire({ icon: 'error', title: 'Error', text: 'Please enter a name' });
                return;
            }
            if (!this.newGroup.handle.trim()) {
                Swal.fire({ icon: 'error', title: 'Error', text: 'Please enter a handle' });
                return;
            }
            this.creating = true;
            try {
                const response = await fetch(urlGroupCreate, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        name: this.newGroup.name.trim(),
                        handle: this.newGroup.handle.trim()
                    })
                });
                const data = await response.json();
                if (data.status === 'success') {
                    Swal.fire({
                        icon: 'success',
                        title: 'Created',
                        text: 'Group created successfully',
                        timer: 1500,
                        showConfirmButton: false
                    });
                    this.showCreateModal = false;
                    this.newGroup = { name: '', handle: '' };
                    this.loadGroups();
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to create group'
                    });
                }
            } catch (error) {
                console.error('Error creating group:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to create group'
                });
            } finally {
                this.creating = false;
            }
        }
    }
};

// Mount the app when DOM is ready. Wrapped in loadVueIfNeeded so
// the app works even if the layout did not load Vue (custom layout).
document.addEventListener('DOMContentLoaded', () => {
    loadVueIfNeeded((err) => {
        if (err) { console.error('Vue load failed:', err); return; }
        const { createApp } = Vue;
        const el = document.getElementById('groups-app');
        const tpl = document.getElementById('groups-app-template');
        if (el && tpl) {
            GroupsApp.template = tpl.innerHTML;
            createApp(GroupsApp).mount('#groups-app');
        }
    });
});
