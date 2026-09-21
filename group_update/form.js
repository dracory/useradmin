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
            groupId: '',
            returnUrl: '',
            errorMessage: '',
            form: {
                status: '',
                name: '',
                handle: '',
                memo: ''
            }
        };
    },
    mounted() {
        this.groupId = GROUP_ID_PLACEHOLDER;
        this.returnUrl = RETURN_URL_PLACEHOLDER;
        this.loadGroup();
    },
    methods: {
        async loadGroup() {
            this.loading = true;
            this.errorMessage = '';
            try {
                const response = await fetch(urlGetGroup, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: new URLSearchParams({
                        action: 'group-fetch-ajax',
                        group_id: this.groupId
                    })
                });
                const result = await response.json();

                if (result.status === 'success') {
                    const d = result.data;
                    this.form.status = d.status || '';
                    this.form.name = d.name || '';
                    this.form.handle = d.handle || '';
                    this.form.memo = d.memo || '';
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to load group', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error loading group:', err);
                Notiflix.Notify.failure('Failed to load group', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                this.loading = false;
            }
        },
        async save(actionType) {
            this.action = actionType;
            this.redirectTo = '';

            if (!this.form.status) {
                Notiflix.Notify.failure('Status is required', { position: 'right-top', timeout: 3000 });
                return;
            }
            if (!this.form.name.trim()) {
                Notiflix.Notify.failure('Name is required', { position: 'right-top', timeout: 3000 });
                return;
            }
            if (!this.form.handle.trim()) {
                Notiflix.Notify.failure('Handle is required', { position: 'right-top', timeout: 3000 });
                return;
            }

            this.saving = true;
            try {
                const response = await fetch(urlUpdateGroup, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        group_id: this.groupId,
                        status: this.form.status,
                        name: this.form.name.trim(),
                        handle: this.form.handle.trim(),
                        memo: this.form.memo.trim(),
                        action: 'group-update-ajax'
                    })
                });
                const result = await response.json();

                if (result.status === 'success') {
                    if (actionType === 'save') {
                        this.redirectTo = this.returnUrl;
                        Notiflix.Notify.success('Group saved successfully', { position: 'right-top', timeout: 3000 });
                        setTimeout(() => {
                            if (this.redirectTo) {
                                window.location.href = this.redirectTo;
                            }
                        }, 3000);
                    } else {
                        Notiflix.Notify.success('Group saved successfully', { position: 'right-top', timeout: 3000 });
                        await this.loadGroup();
                    }
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to save group', { position: 'right-top', timeout: 3000 });
                }
            } catch (err) {
                console.error('Error saving group:', err);
                Notiflix.Notify.failure('Failed to save group', { position: 'right-top', timeout: 3000 });
            } finally {
                this.saving = false;
            }
        }
    }
}).mount('#app-group-update');
});
