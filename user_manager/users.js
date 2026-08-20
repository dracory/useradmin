const UsersApp = {
    data() {
        return {
            // UI state
            loading: true,
            showFilterModal: false,
            showCreateModal: false,
            creating: false,

            // User data
            users: [],
            totalUsers: 0,

            // Pagination
            currentPage: 0,
            perPage: 10,
            jumpToPage: 1,

            // Filters
            filters: {
                status: '',
                first_name: '',
                last_name: '',
                email: '',
                user_id: '',
                created_from: '',
                created_to: ''
            },

            // Sorting
            sortByColumn: 'created_at',
            sortOrder: 'desc',

            // New user
            newUser: {
                first_name: '',
                last_name: '',
                email: ''
            }
        };
    },

    computed: {
        /**
         * Returns total number of pages.
         */
        totalPages() {
            return Math.ceil(this.totalUsers / this.perPage);
        },

        /**
         * Returns a human-readable string describing the current filter state.
         */
        filterStatus() {
            const parts = [];
            if (this.filters.status) parts.push(`status: ${this.filters.status}`);
            if (this.filters.first_name) parts.push(`first name: "${this.filters.first_name}"`);
            if (this.filters.last_name) parts.push(`last name: "${this.filters.last_name}"`);
            if (this.filters.email) parts.push(`email: "${this.filters.email}"`);
            if (this.filters.user_id) parts.push(`user id: "${this.filters.user_id}"`);
            if (this.filters.created_from) parts.push(`from: ${this.filters.created_from}`);
            if (this.filters.created_to) parts.push(`to: ${this.filters.created_to}`);
            
            if (parts.length === 0) return 'Showing all users';
            return 'Showing users with ' + parts.join(', ');
        },

        /**
         * Returns true if any filters are active.
         */
        hasActiveFilters() {
            return this.filters.status !== '' ||
                   this.filters.first_name !== '' ||
                   this.filters.last_name !== '' ||
                   this.filters.email !== '' ||
                   this.filters.user_id !== '' ||
                   this.filters.created_from !== '' ||
                   this.filters.created_to !== '';
        }
    },

    mounted() {
        // Read filters from URL parameters for shareable URLs
        const urlParams = new URLSearchParams(window.location.search);
        
        this.filters.status = urlParams.get('status') || '';
        this.filters.first_name = urlParams.get('first_name') || '';
        this.filters.last_name = urlParams.get('last_name') || '';
        this.filters.email = urlParams.get('email') || '';
        this.filters.user_id = urlParams.get('user_id') || '';
        this.filters.created_from = urlParams.get('created_from') || '';
        this.filters.created_to = urlParams.get('created_to') || '';
        this.sortByColumn = urlParams.get('sort_by') || 'created_at';
        this.sortOrder = urlParams.get('sort_order') || 'desc';
        this.currentPage = parseInt(urlParams.get('page') || '0', 10);
        this.perPage = parseInt(urlParams.get('per_page') || '10', 10);

        this.loadUsers();
    },
    methods: {
        /**
         * Loads users from the server based on current filters, pagination, and sorting.
         */
        async loadUsers() {
            this.loading = true;
            try {
                const response = await fetch(urlUsersLoad, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        page: this.currentPage,
                        per_page: this.perPage,
                        status: this.filters.status,
                        first_name: this.filters.first_name,
                        last_name: this.filters.last_name,
                        email: this.filters.email,
                        user_id: this.filters.user_id,
                        created_from: this.filters.created_from,
                        created_to: this.filters.created_to,
                        sort_by: this.sortByColumn,
                        sort_order: this.sortOrder
                    })
                });
                const data = await response.json();

                if (data.status === 'success') {
                    this.users = data.data?.users || [];
                    this.totalUsers = data.data?.total || 0;
                    // Sync jumpToPage with currentPage to ensure UI consistency
                    this.jumpToPage = this.currentPage + 1;
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to load users'
                    });
                }
            } catch (error) {
                console.error('Error loading users:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'Failed to load users'
                });
            } finally {
                this.loading = false;
            }
        },
        /**
         * Sorts the users by the specified column.
         */
        sortBy(column) {
            if (this.sortByColumn === column) {
                this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
            } else {
                this.sortByColumn = column;
                this.sortOrder = 'asc';
            }
            this.currentPage = 0;
            this.applyFilters();
        },

        /**
         * Navigates to the specified page number.
         */
        goToPage(page) {
            if (page < 0) return;
            const maxPage = Math.ceil(this.totalUsers / this.perPage) - 1;
            if (page > maxPage) page = maxPage;
            this.currentPage = page;
            this.jumpToPage = page + 1;
            this.applyFilters();
        },

        /**
         * Changes the per-page setting and reloads users.
         */
        changePerPage() {
            this.perPage = parseInt(this.perPage, 10);
            this.currentPage = 0;
            this.jumpToPage = 1;
            this.applyFilters();
        },

        /**
         * Opens the filter modal.
         */
        openFilterModal() {
            this.showFilterModal = true;
        },

        /**
         * Closes the filter modal.
         */
        closeFilterModal() {
            this.showFilterModal = false;
        },

        /**
         * Applies the current filters and updates the URL.
         */
        applyFilters() {
            const params = new URLSearchParams();
            if (this.filters.status) params.set('status', this.filters.status);
            if (this.filters.first_name) params.set('first_name', this.filters.first_name);
            if (this.filters.last_name) params.set('last_name', this.filters.last_name);
            if (this.filters.email) params.set('email', this.filters.email);
            if (this.filters.user_id) params.set('user_id', this.filters.user_id);
            if (this.filters.created_from) params.set('created_from', this.filters.created_from);
            if (this.filters.created_to) params.set('created_to', this.filters.created_to);
            params.set('page', this.currentPage);
            params.set('per_page', this.perPage);
            params.set('sort_order', this.sortOrder);
            params.set('sort_by', this.sortByColumn);

            const newUrl = `${window.location.pathname}?${params.toString()}`;
            window.history.pushState({}, '', newUrl);

            this.closeFilterModal();
            this.loadUsers();
        },

        /**
         * Clears all filters and resets to default state.
         */
        clearFilters() {
            this.filters = {
                status: '',
                first_name: '',
                last_name: '',
                email: '',
                user_id: '',
                created_from: '',
                created_to: ''
            };
            this.currentPage = 0;
            this.applyFilters();
        },

        /**
         * Returns the sort icon class for a column.
         */
        sortIcon(column) {
            if (this.sortByColumn !== column) return 'bi bi-arrow-down-up text-muted';
            return this.sortOrder === 'asc' ? 'bi bi-arrow-up' : 'bi bi-arrow-down';
        },
        statusClass(status) {
            switch (status) {
                case 'active': return 'bg-success';
                case 'inactive': return 'bg-danger';
                case 'unverified': return 'bg-info';
                case 'deleted': return 'bg-secondary';
                default: return 'bg-light text-dark';
            }
        },
        /**
         * Deletes the specified user after confirmation.
         */
        async deleteUser(user) {
            const result = await Swal.fire({
                icon: 'warning',
                title: 'Delete User?',
                text: `Are you sure you want to delete ${user.first_name} ${user.last_name}?`,
                showCancelButton: true,
                confirmButtonText: 'Delete',
                cancelButtonText: 'Cancel',
                confirmButtonColor: '#dc3545'
            });

            if (!result.isConfirmed) return;

            try {
                const response = await fetch(urlUserDelete, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ user_id: user.id })
                });

                const data = await response.json();

                if (data.status === 'success') {
                    Swal.fire({
                        icon: 'success',
                        title: 'Deleted',
                        text: 'User deleted successfully',
                        timer: 1500,
                        showConfirmButton: false
                    });
                    this.loadUsers();
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to delete user'
                    });
                }
            } catch (error) {
                console.error('Error deleting user:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to delete user'
                });
            }
        },
        /**
         * Creates a new user.
         */
        async createUser() {
            if (!this.newUser.first_name || !this.newUser.last_name || !this.newUser.email) {
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Please fill in all fields'
                });
                return;
            }
            this.creating = true;
            try {
                const response = await fetch(urlUserCreate, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.newUser)
                });
                const data = await response.json();
                if (data.status === 'success') {
                    Swal.fire({
                        icon: 'success',
                        title: 'Created',
                        text: 'User created successfully',
                        timer: 1500,
                        showConfirmButton: false
                    });
                    this.showCreateModal = false;
                    this.newUser = { first_name: '', last_name: '', email: '' };
                    this.currentPage = 0;
                    this.loadUsers();
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to create user'
                    });
                }
            } catch (error) {
                console.error('Error creating user:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to create user'
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
        const el = document.getElementById('users-app');
        if (el) {
            createApp(UsersApp).mount('#users-app');
        }
    });
});
