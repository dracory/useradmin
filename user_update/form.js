// Wrapped in loadVueIfNeeded so the app works even if the layout
// did not load Vue (custom layout).
loadVueIfNeeded((err) => {
    if (err) { console.error('Vue load failed:', err); return; }
    const { createApp } = Vue;

    createApp({
    template: document.getElementById('app-user-update-template').innerHTML,
    data() {
        return {
            loading: true,
            saving: false,
            savingSection: '',
            passwordSaving: false,
            passwordCollapsed: false,
            notesCollapsed: false,
            locationCollapsed: false,
            detailsCollapsed: false,
            statusCollapsed: false,
            rolesCollapsed: false,
            showPassword: false,
            showPasswordConfirm: false,
            passwordForm: {
                password: '',
                password_confirm: ''
            },
            action: '',
            redirectTo: '',
            userId: '',
            returnUrl: '',
            form: {
                status: '',
                first_name: '',
                last_name: '',
                email: '',
                business_name: '',
                phone: '',
                country: '',
                timezone: '',
                memo: '',
                role: '',
                role_ids: [],
                group_ids: []
            },
            originalEmail: '',
            countries: [],
            timezones: [],
            roles: [],
            groups: [],
            fieldStatus: {
                first_name: true,
                last_name: true,
                email: true,
                business_name: true,
                phone: true
            }
        };
    },
    computed: {
        hasUnreadableFields() {
            return Object.values(this.fieldStatus).some(v => !v);
        },
        warningFields() {
            const fields = [];
            if (!this.fieldStatus.first_name) fields.push('first name');
            if (!this.fieldStatus.last_name) fields.push('last name');
            if (!this.fieldStatus.email) fields.push('email');
            if (!this.fieldStatus.business_name) fields.push('business name');
            if (!this.fieldStatus.phone) fields.push('phone');
            return fields;
        }
    },
    mounted() {
        this.userId = USER_ID_PLACEHOLDER;
        this.returnUrl = RETURN_URL_PLACEHOLDER;
        this.loadUser(true);
    },
    methods: {
        async loadUser(showSpinner) {
            if (showSpinner) {
                this.loading = true;
            }
            this.errorMessage = '';
            try {
                const response = await fetch(urlGetUser, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: new URLSearchParams({
                        action: 'user-fetch-ajax',
                        user_id: this.userId
                    })
                });
                const result = await response.json();

                if (result.status === 'success') {
                    const d = result.data;
                    this.form.status = d.status || '';
                    this.form.first_name = d.first_name || '';
                    this.form.last_name = d.last_name || '';
                    this.form.email = d.email || '';
                    this.form.business_name = d.business_name || '';
                    this.form.phone = d.phone || '';
                    this.form.country = d.country || '';
                    this.form.timezone = d.timezone || '';
                    this.form.memo = d.memo || '';
                    this.form.role = d.role || '';
                    this.originalEmail = d.email || '';
                    this.countries = d.countries || [];
                    this.timezones = d.timezones || [];
                    this.roles = d.roles || [];
                    this.groups = d.groups || [];
                    this.form.role_ids = d.user_role_ids || [];
                    this.form.group_ids = d.user_group_ids || [];
                    if (d.field_status) {
                        this.fieldStatus = { ...this.fieldStatus, ...d.field_status };
                    }
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to load user', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error loading user:', err);
                Notiflix.Notify.failure('Failed to load user', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                if (showSpinner) {
                    this.loading = false;
                }
            }
        },
        async onCountryChange() {
            this.form.timezone = '';
            if (!this.form.country) {
                this.timezones = [];
                return;
            }
            try {
                const response = await fetch(urlGetTimezones, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: new URLSearchParams({
                        action: 'get-timezones-ajax',
                        country_code: this.form.country
                    })
                });
                const result = await response.json();
                if (result.status === 'success') {
                    this.timezones = result.data.timezones || [];
                }
            } catch (err) {
                console.error('Error loading timezones:', err);
            }
        },
        async changePassword() {
            if (this.passwordForm.password.length < 8) {
                Notiflix.Notify.failure('Password must be at least 8 characters', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (this.passwordForm.password !== this.passwordForm.password_confirm) {
                Notiflix.Notify.failure('Passwords do not match', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }

            this.passwordSaving = true;
            try {
                const response = await fetch(urlUpdatePassword, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        action: 'password-update-ajax',
                        user_id: this.userId,
                        password: this.passwordForm.password
                    })
                });
                const result = await response.json();

                if (result.status === 'success') {
                    this.passwordForm.password = '';
                    this.passwordForm.password_confirm = '';
                    this.showPassword = false;
                    this.showPasswordConfirm = false;
                    Notiflix.Notify.success('Password changed successfully', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to change password', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error changing password:', err);
                Notiflix.Notify.failure('Failed to change password', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                this.passwordSaving = false;
            }
        },
        sectionConfig() {
            return {
                status: {
                    url: urlUpdateStatus,
                    action: 'status-update-ajax',
                    validate: () => this.form.status ? '' : 'Status is required',
                    body: () => ({
                        status: this.form.status.trim()
                    })
                },
                details: {
                    url: urlUpdateDetails,
                    action: 'details-update-ajax',
                    validate: () => {
                        if (!this.form.email.trim()) return 'Email is required';
                        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                        if (!emailRegex.test(this.form.email.trim())) return 'Invalid email address';
                        return '';
                    },
                    body: () => ({
                        first_name: this.form.first_name.trim(),
                        last_name: this.form.last_name.trim(),
                        email: this.form.email.trim(),
                        business_name: this.form.business_name.trim(),
                        phone: this.form.phone.trim()
                    })
                },
                roles: {
                    url: urlUpdateRoles,
                    action: 'roles-update-ajax',
                    validate: () => '',
                    body: () => ({
                        role: this.form.role.trim(),
                        role_ids: this.form.role_ids,
                        group_ids: this.form.group_ids
                    })
                },
                location: {
                    url: urlUpdateLocation,
                    action: 'location-update-ajax',
                    validate: () => '',
                    body: () => ({
                        country: this.form.country,
                        timezone: this.form.timezone
                    })
                },
                notes: {
                    url: urlUpdateNotes,
                    action: 'notes-update-ajax',
                    validate: () => '',
                    body: () => ({
                        memo: this.form.memo.trim()
                    })
                }
            };
        },
        async postSection(section) {
            const conf = this.sectionConfig()[section];
            const validationError = conf.validate();
            if (validationError) {
                Notiflix.Notify.failure(validationError, {
                    position: 'right-top',
                    timeout: 3000,
                });
                return false;
            }
            const response = await fetch(conf.url, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    ...conf.body(),
                    user_id: this.userId,
                    action: conf.action
                })
            });
            const result = await response.json();
            if (result.status === 'success') {
                return true;
            }
            Notiflix.Notify.failure(result.message || 'Failed to save', {
                position: 'right-top',
                timeout: 3000,
            });
            return false;
        },
        async saveSection(section) {
            this.action = section;
            this.savingSection = section;
            try {
                if (await this.postSection(section)) {
                    await this.loadUser(false);
                    Notiflix.Notify.success('Saved successfully', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error saving ' + section + ':', err);
                Notiflix.Notify.failure('Failed to save', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                this.savingSection = '';
            }
        }
    }
}).mount('#app-user-update');
});
