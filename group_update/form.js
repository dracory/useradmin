// Wrapped in loadVueIfNeeded so the app works even if the layout
// did not load Vue (custom layout).
loadVueIfNeeded((err) => {
    if (err) { console.error('Vue load failed:', err); return; }
    const { createApp } = Vue;

    createApp({
    template: document.getElementById('app-group-update-template').innerHTML,
    data() {
        return {
            loading: true,
            savingSection: '',
            statusCollapsed: false,
            detailsCollapsed: false,
            notesCollapsed: false,
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
        this.loadGroup(true);
    },
    methods: {
        async loadGroup(showSpinner) {
            if (showSpinner) {
                this.loading = true;
            }
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
                if (showSpinner) {
                    this.loading = false;
                }
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
                        if (!this.form.name.trim()) return 'Name is required';
                        if (!this.form.handle.trim()) return 'Handle is required';
                        return '';
                    },
                    body: () => ({
                        name: this.form.name.trim(),
                        handle: this.form.handle.trim()
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
                    group_id: this.groupId,
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
            this.savingSection = section;
            try {
                if (await this.postSection(section)) {
                    await this.loadGroup(false);
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
}).mount('#app-group-update');
});
