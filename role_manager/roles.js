const RolesApp = {
    data() {
        return {
            loading: true,
            showCreateModal: false,
            creating: false,
            roles: [],
            urlRoleUpdate: urlRoleUpdate,
            newRole: {
                name: '',
                handle: ''
            }
        };
    },
    mounted() {
        this.loadRoles();
    },
    methods: {
        async loadRoles() {
            this.loading = true;
            try {
                const response = await fetch(urlRolesLoad, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({})
                });
                const data = await response.json();

                if (data.status === 'success') {
                    this.roles = data.data?.roles || [];
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to load roles'
                    });
                }
            } catch (error) {
                console.error('Error loading roles:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'Failed to load roles'
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
        async deleteRole(role) {
            const result = await Swal.fire({
                icon: 'warning',
                title: 'Delete Role?',
                text: `Are you sure you want to delete ${role.name}?`,
                showCancelButton: true,
                confirmButtonText: 'Delete',
                cancelButtonText: 'Cancel',
                confirmButtonColor: '#dc3545'
            });

            if (!result.isConfirmed) return;

            try {
                const response = await fetch(urlRoleDelete, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ role_id: role.id })
                });

                const data = await response.json();

                if (data.status === 'success') {
                    Swal.fire({
                        icon: 'success',
                        title: 'Deleted',
                        text: 'Role deleted successfully',
                        timer: 1500,
                        showConfirmButton: false
                    });
                    this.loadRoles();
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to delete role'
                    });
                }
            } catch (error) {
                console.error('Error deleting role:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to delete role'
                });
            }
        },
        async createRole() {
            if (!this.newRole.name.trim()) {
                Swal.fire({ icon: 'error', title: 'Error', text: 'Please enter a name' });
                return;
            }
            if (!this.newRole.handle.trim()) {
                Swal.fire({ icon: 'error', title: 'Error', text: 'Please enter a handle' });
                return;
            }
            this.creating = true;
            try {
                const response = await fetch(urlRoleCreate, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        name: this.newRole.name.trim(),
                        handle: this.newRole.handle.trim()
                    })
                });
                const data = await response.json();
                if (data.status === 'success') {
                    Swal.fire({
                        icon: 'success',
                        title: 'Created',
                        text: 'Role created successfully',
                        timer: 1500,
                        showConfirmButton: false
                    });
                    this.showCreateModal = false;
                    this.newRole = { name: '', handle: '' };
                    this.loadRoles();
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to create role'
                    });
                }
            } catch (error) {
                console.error('Error creating role:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to create role'
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
        const el = document.getElementById('roles-app');
        const tpl = document.getElementById('roles-app-template');
        if (el && tpl) {
            RolesApp.template = tpl.innerHTML;
            createApp(RolesApp).mount('#roles-app');
        }
    });
});
