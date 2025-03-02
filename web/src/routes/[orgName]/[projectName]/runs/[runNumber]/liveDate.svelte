<script lang="ts">
	import Dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';

	Dayjs.extend(duration);

	interface LiveDateProps {
		startedAt: string;
		finishedAt?: string;
	}

	let time = $state('');

	const { startedAt, finishedAt }: LiveDateProps = $props();

	$effect(() => {
		const calculateDuration = () => {
			const d = Dayjs.duration(Dayjs(finishedAt ?? Date.now()).diff(Dayjs(startedAt)));

			const days = d.days() ? `${d.days()}d ` : '';
			const hours = d.hours() ? `${d.hours()}h ` : '';
			const minutes = d.minutes() ? `${d.minutes()}m ` : '';
			const seconds = `${d.seconds()}s`;

			time = `${days}${hours}${minutes}${seconds}`.trim();

			if (!finishedAt)
				setTimeout(() => {
					calculateDuration();
				}, 100);
		};

		calculateDuration();
	});
</script>

{time}
