// Wrapped in loadVueIfNeeded so the app works even if the layout
// did not load Vue (custom layout).
loadVueIfNeeded((err) => {
    if (err) { console.error('Vue load failed:', err); return; }
    const { createApp } = Vue;

    createApp({
    data() {
        return {
            loading: true,
            saving: false,
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
                role: ''
            },
            originalEmail: '',
            countries: [],
            timezones: [],
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
        this.loadUser();
    },
    methods: {
        async loadUser() {
            this.loading = true;
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
                this.loading = false;
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
        async save(actionType) {
            this.action = actionType;
            this.redirectTo = '';

            if (!this.form.status) {
                Notiflix.Notify.failure('Status is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (!this.form.first_name.trim()) {
                Notiflix.Notify.failure('First name is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (!this.form.last_name.trim()) {
                Notiflix.Notify.failure('Last name is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (!this.form.email.trim()) {
                Notiflix.Notify.failure('Email is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(this.form.email.trim())) {
                Notiflix.Notify.failure('Invalid email address', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (!this.form.country) {
                Notiflix.Notify.failure('Country is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (!this.form.timezone) {
                Notiflix.Notify.failure('Timezone is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }

            this.saving = true;
            try {
                const payload = {
                    user_id: this.userId,
                    status: this.form.status.trim(),
                    first_name: this.form.first_name.trim(),
                    last_name: this.form.last_name.trim(),
                    email: this.form.email.trim(),
                    business_name: this.form.business_name.trim(),
                    phone: this.form.phone.trim(),
                    country: this.form.country,
                    timezone: this.form.timezone,
                    memo: this.form.memo.trim(),
                    role: this.form.role.trim()
                };

                const response = await fetch(urlUpdateUser, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        ...payload,
                        action: 'user-update-ajax'
                    })
                });
                const result = await response.json();

                if (result.status === 'success') {
                    if (actionType === 'save') {
                        this.redirectTo = this.returnUrl;
                        Notiflix.Notify.success('User saved successfully', {
                            position: 'right-top',
                            timeout: 3000,
                        });
                        setTimeout(() => {
                            if (this.redirectTo) {
                                window.location.href = this.redirectTo;
                            }
                        }, 3000);
                    } else {
                        Notiflix.Notify.success('User saved successfully', {
                            position: 'right-top',
                            timeout: 3000,
                        });
                        // Reload user data after successful save
                        await this.loadUser();
                    }
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to save user', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error saving user:', err);
                Notiflix.Notify.failure('Failed to save user', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                this.saving = false;
            }
        }
    }
}).mount('#app-user-update');
});
